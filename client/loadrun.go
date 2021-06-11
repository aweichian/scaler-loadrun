package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func HttpClient(number int) {
	t := time.Now()
	fmt.Printf("do http request: %v, number: %d\n", t, number)
}

func LoadRun(p_url string) (err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, p_url, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout: 5 * time.Second,
				}
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}

	var response *http.Response
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return
	}

	var respBytes []byte
	respBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return
	}

	var result interface{}
	if response.StatusCode == 200 {
		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return
		}
	}
	//fmt.Printf("result: %v\n",result)
	return nil
}

func main() {
	var url string
	var pool int
	flag.StringVar(&url, "url", "localhost:8080", "load run url")
	flag.IntVar(&pool, "pool", 100, "load run pool")
	flag.Parse()
	ch := make(chan struct{}, pool)
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- struct{}{}
			time.Sleep(time.Millisecond * 500)
		}()
	}

	//for i := 0; i < pool; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		if ch != nil {
	//			<-ch
	//			LoadRun(url)
	//		}
	//	}()
	//}
	//
	//
	//wg.Wait()
	//for {
	//  select {
	//  case <-ch:
	//      count ++
	//      //HttpClient(count)
	//      LoadRun("http://localhost:8080/v1/task?type=cpu&number=10000000")
	//  default:
	//  }
	//}
}
