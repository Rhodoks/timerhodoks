package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var info string = `{\"Method\":\"GET\",\"Url\":\"http://www.baidu.com\",\"Headers\":\"{\\\"Hello\\\":\\\"World\\\"}\",\"Body\":\"\",\"Timeout\":\"10s\",\"ExpectCode\":200}`

func SendJob(data string) {
	targetUrl := "http://localhost:18001/api/job"

	payload := strings.NewReader(data)

	req, _ := http.NewRequest("PUT", targetUrl, payload)

	// fmt.Println(err)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	fmt.Println(err, res)
}

func main() {
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond * 10)
		SendJob(`{"Name":"Test","CronLine":"* * * * * * *","ExecutorType":"Http","ExecutorInfo":"` + info + `","RetryNum":1}`)
	}
}
