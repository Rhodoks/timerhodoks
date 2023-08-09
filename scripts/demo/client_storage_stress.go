package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

func SendJob(targetUrl string, data string) {
	payload := strings.NewReader(data)
	req, _ := http.NewRequest("PUT", targetUrl, payload)
	req.Header.Add("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}

var targetUrl = []string{"http://127.0.0.1:18001/api/job", "http://127.0.0.1:28001/api/job", "http://127.0.0.1:38001/api/job"}
var jobinfo = `{` +
	`\"Method\":\"GET\",` +
	`\"Url\":\"http://127.0.0.1:8000/stats\",` +
	`\"Headers\":\"{}\",` +
	`\"Body\":\"\",` +
	`\"Timeout\":\"1s\",` +
	`\"ExpectCode\":200` +
	`}`

func RandomlySendJob(cronline string, Cnt int) {
	url := targetUrl[rand.Int31()%3]
	SendJob(url, fmt.Sprintf(`{"Name":"Test","CronLine":"`+cronline+`","ExecutorType":"Http","ExecutorInfo":"`+jobinfo+`","RetryNum":1,"Cnt":%d}`, Cnt))
}

func Run(chanF chan struct{}) {
}

func main() {
	every_minute := flag.Int("every_minute", 0, "insert 60 jobs each triggers every minute")
	dead_job := flag.Int("dead_job", 0, "insert 60 jobs each triggers every minute")
	flag.Parse()
	fmt.Printf("Insert %d * %d = %d jobs, should trigger exactly %d times per in second each minute.", *every_minute, 60, *every_minute*60, *every_minute)
	for i := 0; i < 60; i++ {
		RandomlySendJob(fmt.Sprintf("%d * * * * * *", i), *every_minute)
	}
	fmt.Printf("Insert %d jobs, which will not trigger for now.", *dead_job)
	RandomlySendJob(fmt.Sprintf("* * * * * * 2024"), *dead_job)
}
