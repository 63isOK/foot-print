package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var dataChannel = make(chan interface{})

func main() {
	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			Receive()
			// time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("协程数:", runtime.NumGoroutine())
		}
	}()

	Start()
}

// Start 启动服务
func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				w.Write(getFailedResponse("请求参数错误"))
			}

			if err != nil {
				w.Write(getFailedResponse("请求参数错误"))
				fmt.Println(err)
			} else {
				w.Write(getSuccessResponse())
			}
		}()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		r.Body.Close()

		// 数据校验
		var postData PostData
		err = json.Unmarshal(body, &postData)
		if err == nil {
			dataChannel <- postData
			return
		}
		var instanceConfigs InstanceConfigs
		err = json.Unmarshal(body, &instanceConfigs)
		if err == nil {
			dataChannel <- instanceConfigs
			return
		}
	})
	handler := http.TimeoutHandler(http.DefaultServeMux,
		time.Duration(3)*time.Second,
		string(getFailedResponse("请求超时")))

	http.ListenAndServe(":5009", handler)
}

func Receive() {
	data := <-dataChannel
	switch data := data.(type) {
	case InstanceConfigs:
		fmt.Println(data)
	case PostData:
		fmt.Println(data)
	}
}

func getFailedResponse(msg string) []byte {
	resp := CommonResult{
		Success: false,
		Message: msg,
	}

	data, _ := json.Marshal(resp)
	return data
}

func getSuccessResponse() []byte {
	resp := CommonResult{
		Success: true,
		Message: "",
	}

	data, _ := json.Marshal(resp)
	return data
}
