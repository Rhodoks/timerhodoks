package main

import (
	"flag"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func SendJob(targetUrl string, data string) {
	payload := strings.NewReader(data)
	req, _ := http.NewRequest("PUT", targetUrl, payload)
	req.Header.Add("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}

var targetUrl = []string{"http://localhost:18001/api/job", "http://localhost:28001/api/job", "http://localhost:38001/api/job"}
var jobinfo = `{` +
	`\"Method\":\"GET\",` +
	`\"Url\":\"http://www.baidu.com\",` +
	`\"Headers\":\"{\\\"Hello\\\":\\\"World\\\"}\",` +
	`\"Body\":\"\",` +
	`\"Timeout\":\"1s\",` +
	`\"ExpectCode\":200` +
	`}`

func SendRandJob() {
	url := targetUrl[rand.Int31()%3]
	SendJob(url, `{"Name":"Test","CronLine":"* * * * * * *","ExecutorType":"Http","ExecutorInfo":"`+jobinfo+`","RetryNum":1}`)
}

func Run(chanF chan struct{}) {
}

func main() {
	num := flag.Int("num", 1, "number of job")
	flag.Parse()
	for i := 0; i < *num; i++ {
		SendRandJob()
		time.Sleep(10 * time.Millisecond)
	}
}
