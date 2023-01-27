package ratelimit

import (
	"encoding/binary"
	"github.com/cespare/xxhash/v2"
	"github.com/ego008/goutils/json"
	"sync"
	"sync/atomic"
	"time"
)

// data struct
// 	bit := []byte{
//		1, 2, 3, 4,
//		1, 2, 3, 4,
//		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
//		1, 2, 3, 4}
//	fmt.Println("bit len", len(bit)) // len 36
//	fmt.Println(bit[:4]) // latest hour
//	fmt.Println(bit[4:8]) // latest quarter
//	fmt.Println(bit[8:32]) // each hour, new in front
//	fmt.Println(bit[32:]) // each quarter, new in front

//  Usage:
// maxEntry, rateLimitDay, rateLimitHour := 1000, 2000, 100
// limiter = NewCache(maxEntry, rateLimitDay, rateLimitHour)
// userIp := getUserIpFromReq(httpRequest)
// cntDay, cntHour, underRateLimit := limiter.Incr(uint64(time.Now().UTC().Unix()), userIp)
//
//	if !underRateLimit {
//	    log.Printf("%s over rate limit, current D[%d] H[%d] \n", userIp, cntDay, cntHour)
//	}
//
// // tools
// byteData := limiter.Dump()
// // save byteData in db or disk
// limiter.Load(byteData) // load byteData to limiter
// // do every n hour
// limiter.RemoveExpiredEntry()

const (
	hourSeconds    = 3600 // 60*60
	quarterSeconds = 900  // 60*15
)

type Em map[uint64]*entry // key: xxhash.Sum64String(ip)
type Cache struct {
	MaxEntries int // not used
	dayLimit   int
	hourLimit  int
	cache      atomic.Value
	lock       sync.RWMutex
}

type entry struct {
	Key  uint64 // xxhash.Sum64String(ip)
	Data []byte // make([]byte, 36, 36)
}

func (e *entry) genByte(curStamp uint64) (curDayValue, curHourValue uint32) {
	//curHour := uint64(math.Ceil(float64(curStamp) / hourSeconds))
	//curQuarter := uint64(math.Ceil(float64(curStamp) / quarterSeconds))

	curHour := curStamp / hourSeconds
	curQuarter := curStamp / quarterSeconds

	// old
	latestHour := b2i32(e.Data[:4])
	latestQuarter := b2i32(e.Data[4:8])

	diffQ := curQuarter - uint64(latestQuarter)
	qLst := make([]byte, 4, 4)
	switch {
	case diffQ == 0:
		qLst = e.Data[32:]
		if qLst[0] < 254 {
			qLst[0]++
		}
	case diffQ < 4:
		qLst[0] = 1
		ito := 4 - diffQ
		for i := uint64(0); i < ito; i++ {
			qLst[diffQ+i] = e.Data[32+i]
		}
	case diffQ >= 4:
		qLst[0] = 1
	}

	curHourValue = uint32(qLst[0] + qLst[1] + qLst[2] + qLst[3])
	var curHourValue8 uint8
	if curHourValue >= 255 {
		curHourValue8 = 255
	} else {
		curHourValue8 = uint8(curHourValue)
	}

	diffH := curHour - uint64(latestHour)
	hLst := make([]byte, 24, 24)
	switch {
	case diffH == 0:
		hLst = e.Data[8:32]
		hLst[0] = curHourValue8
	case diffH < 24:
		hLst[0] = curHourValue8
		ito := 24 - diffH
		for i := uint64(0); i < ito; i++ {
			hLst[diffH+i] = e.Data[8+i]
		}
	case diffH >= 24:
		hLst[0] = curHourValue8
	}

	for i := 0; i < 24; i++ {
		curDayValue += uint32(hLst[i])
	}

	var curHourB, curQuarterB []byte
	if diffH > 0 {
		curHourB = i2b32(uint32(curHour))
	} else {
		curHourB = e.Data[:4]
	}
	if diffQ > 0 {
		curQuarterB = i2b32(uint32(curQuarter))
	} else {
		curQuarterB = e.Data[4:8]
	}

	e.Data = e.Data[:0]
	e.Data = append(e.Data, curHourB...)
	e.Data = append(e.Data, curQuarterB...)
	e.Data = append(e.Data, hLst...)
	e.Data = append(e.Data, qLst...)

	return
}

func NewCache(maxEntry, dayLimit, hourLimit int) *Cache {
	// limited: max 255 per quarter, hourMaxLimit: 255*4, dayMaxLimit: 255*24
	hourMaxLimit := 255 * 4
	dayMaxLimit := 255 * 24
	if dayLimit < 0 {
		dayLimit = 0
	} else if dayLimit > dayMaxLimit {
		dayLimit = dayMaxLimit
	}
	if hourLimit < 0 {
		hourLimit = 0
	} else if hourLimit > hourMaxLimit {
		hourLimit = hourMaxLimit
	}
	c := &Cache{
		MaxEntries: maxEntry, // not used
		dayLimit:   dayLimit,
		hourLimit:  hourLimit,
	}
	c.cache.Store(Em{})
	return c
}

func (c *Cache) Incr(curStamp uint64, key string) (uint32, uint32, bool) {
	if c.hourLimit == 0 && c.dayLimit == 0 {
		return 0, 0, false
	}

	underRateLimit := true
	keyI64 := xxhash.Sum64String(key)

	c.lock.Lock()
	defer c.lock.Unlock()

	cache := c.cache.Load().(Em)
	if ee, ok := cache[keyI64]; ok {
		dayValue, hourValue := ee.genByte(curStamp)
		if int(dayValue) > c.dayLimit || int(hourValue) > c.hourLimit {
			underRateLimit = false
		}
		cache[keyI64] = ee
		c.cache.Store(cache)
		return dayValue, hourValue, underRateLimit
	}

	// new item
	item := &entry{
		Key:  keyI64,
		Data: make([]byte, 36, 36),
	}
	item.genByte(curStamp)
	cache[keyI64] = item
	c.cache.Store(cache)

	return 1, 1, underRateLimit
}

// Load from []byte
func (c *Cache) Load(jb []byte) {
	var items []entry
	err := json.Unmarshal(jb, &items)
	if err != nil {
		return
	}
	curHour := uint32(time.Now().UTC().Unix() / hourSeconds)
	cache := c.cache.Load().(Em)
	for _, v := range items {
		if (curHour - b2i32(v.Data[:4])) > 24 {
			continue
		}
		cache[v.Key] = &v
	}
	c.cache.Store(cache)
}

// Dump to []byte
func (c *Cache) Dump() []byte {
	cache := c.cache.Load().(Em)
	if len(cache) == 0 {
		return []byte{}
	}
	curHour := uint32(time.Now().UTC().Unix() / hourSeconds)
	var items []entry
	for _, v := range cache {
		if (curHour - b2i32(v.Data[:4])) > 24 {
			continue
		}
		items = append(items, *v)
	}
	jb, _ := json.Marshal(items)
	return jb
}

func (c *Cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	cache := c.cache.Load().(Em)
	return len(cache)
}

func (c *Cache) RemoveExpiredEntry() {
	curHour := uint32(time.Now().UTC().Unix() / hourSeconds)
	c.lock.Lock()
	defer c.lock.Unlock()
	cache := c.cache.Load().(Em)
	for k, v := range cache {
		if (curHour - b2i32(v.Data[:4])) > 24 {
			delete(cache, k)
		}
	}
	c.cache.Store(cache)
}

// utils

func i2b32(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
func b2i32(v []byte) uint32 {
	return binary.BigEndian.Uint32(v)
}
