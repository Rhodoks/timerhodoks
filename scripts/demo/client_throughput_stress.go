package main

import (
	"flag"
	"fmt"
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

func main() {
	num := flag.Int("num", 1, "number of job")
	flag.Parse()
	fmt.Printf("Insert %d jobs, each should trigger per second.", *num)
	SendJob(targetUrl[0], fmt.Sprintf(`{"Name":"Test","CronLine":"* * * * * * *","ExecutorType":"Http","ExecutorInfo":"`+jobinfo+`","RetryNum":1,"Cnt":%d}`, *num))
}
