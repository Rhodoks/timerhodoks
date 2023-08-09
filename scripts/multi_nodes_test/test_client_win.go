package main

import (
	"fmt"
	"net/http"
	"strings"
)

func SendJob(targetUrl string, data string) {

	payload := strings.NewReader(data)

	req, _ := http.NewRequest("PUT", targetUrl, payload)

	// fmt.Println(err)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	fmt.Println(err, res)
}

func main() {
	SendJob("http://localhost:18001/job", `{"Name":"Test","CronLine":"0,10,20,30,40,50 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 1 >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	SendJob("http://localhost:28001/job", `{"Name":"Test","CronLine":"5,15,25,35,45,55 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 2 >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	SendJob("http://localhost:38001/job", `{"Name":"Test","CronLine":"1,2,3,4,5,6,7,8,9 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 3 >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	// for i := 0; i < 10000; i++ {
	// 	SendJob(`{"Name":"Test","CronLine":"0,10,20,30,40,50 * * * * * 2030","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 1 >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	// }
}
