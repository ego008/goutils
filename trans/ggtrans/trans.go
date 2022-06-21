package ggtrans

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/ego008/goutils/json"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

// Translate Translate("127.0.0.1:1080", "userAgent", "你好", "zh-CN", "en") // auto
func Translate(curSock5, ua, source, sourceLang, targetLang string) (string, error) {
	var tr = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 nil,
	}
	if len(curSock5) > 0 {
		dialer, err := proxy.SOCKS5("tcp", curSock5,
			nil,
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		)
		if err == nil {
			tr.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, e := dialer.Dial(network, addr)
				return c, e
			}
		} else {
			fmt.Println("DialContext err", err)
		}
	}

	var HttpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}

	var text []string
	var result []interface{}

	encodedSource := encodeURIComponent(source)
	bsUrl := "https://translate.googleapis.com/translate_a/single"
	tranUrl := bsUrl + "?client=gtx&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	client := HttpClient
	req, _ := http.NewRequest("GET", tranUrl, nil)
	if len(ua) > 0 {
		req.Header.Set("User-Agent", ua)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response body ")
	}

	bReq := strings.Contains(string(body), `<title>Error 400 (Bad Request)`)
	if bReq {
		return "", errors.New("Error 400 (Bad Request) ")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", errors.New("Error Unmarshal data ")
	}

	if len(result) > 0 {
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		cText := strings.Join(text, "")

		return cText, nil
	} else {
		return "", errors.New("No translated data in response ")
	}
}
