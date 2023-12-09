package racecondition

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type RequestCrack struct {
	Proxy   *url.URL
	Url     string
	Request *http.Request
	Body    bytes.Buffer
}

var startChan = make(chan bool, 0)

func (r RequestCrack) sendRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if r.Proxy != nil {
		tr.Proxy = http.ProxyURL(r.Proxy)
	}

	payload := ioutil.NopCloser(&r.Body)
	req, err := http.NewRequest(r.Request.Method, r.Url, payload)
	if err != nil {
		log.Println(err)
		return
	}

	for key, values := range r.Request.Header {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	client := &http.Client{
		Transport: tr,
	}

	timeSend := time.Now()
	<-startChan
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("[WORKER %d] Time: %v\nStatus: %v\nBody: %s\n-----------------------------------------------------\n",
		id, timeSend.Format("15:04:05.000 02-01-2006"), resp.Status, string(body))
}

func Start() {
	if requestFile == "" {
		fmt.Println("Please specify all arguments!")
		flagCmd.PrintDefaults()
		return
	}

	var proxyUrl *url.URL
	var err error

	if len(proxy) > 0 {
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			panic(err)
		}
	}

	urlStr, httpRequest, body := parseReqFromFile(host, requestFile, https)
	if len(host) < 1 {
		host = httpRequest.Host
	}

	client := RequestCrack{
		Proxy:   proxyUrl,
		Url:     urlStr,
		Request: httpRequest,
		Body:    body,
	}

	totalWorker := numWorker
	var wg1 sync.WaitGroup
	if len(requestFile1) > 0 {
		urlStr1, httpRequest1, body1 := parseReqFromFile(host, requestFile1, https)
		client1 := RequestCrack{
			Proxy:   proxyUrl,
			Url:     urlStr1,
			Request: httpRequest1,
			Body:    body1,
		}

		for i := 0; i < numWorker; i++ {
			wg1.Add(1)
			go func(i int) {
				client1.sendRequest(&wg1, i+1)
			}(i)
		}
		totalWorker += numWorker
	}

	var wg sync.WaitGroup
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go func(i int) {
			client.sendRequest(&wg, i+1)
		}(i)
	}

	for i := 0; i < totalWorker; i++ {
		startChan <- true
	}

	wg.Wait()
	wg1.Wait()
}

func parseReqFromFile(host string, requestFile string, https bool) (urlStr string, httpRequest *http.Request, body bytes.Buffer) {
	content, err := ioutil.ReadFile(requestFile)
	if err != nil {
		panic(err)
	}

	httpRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(content)))
	if err != nil && err != io.ErrUnexpectedEOF {
		if strings.Contains(err.Error(), `HTTP version "HTTP/2"`) {
			fmt.Println(`To fix it: Replace "HTTP/2" to "HTTP/2.0" in request file`)
		}
		fmt.Println("parse http request from file error", err)
		return
	}

	if len(host) < 1 {
		host = httpRequest.Host
	}

	urlStr = "http://" + host + httpRequest.RequestURI
	if https {
		urlStr = "https://" + host + httpRequest.RequestURI
	}

	bodyLen, err := body.ReadFrom(httpRequest.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			fmt.Println("read body request error", err)
			return
		}
		httpRequest.Header.Set("Content-Length", fmt.Sprintf("%d", bodyLen))
	}

	return
}
