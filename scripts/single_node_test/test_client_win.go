package main

import (
	"fmt"
	"net/http"
	"strings"
)

func SendJob(data string) {
	targetUrl := "http://localhost:18001/api/job"

	payload := strings.NewReader(data)

	req, _ := http.NewRequest("PUT", targetUrl, payload)

	// fmt.Println(err)

	req.Header.Add("Content-Type", "application/json")

	err, res := http.DefaultClient.Do(req)

	fmt.Println(err, res)
}

// {"Shell":"PowerShell","Command":"echo 1 $JOBID $TRIGGERTIME >> test.txt","Timeout":"10s"}
func main() {
	SendJob(`{"Name":"Test","CronLine":"0,10,20,30,40,50 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 1 $JOBID $TRIGGERTIME >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	SendJob(`{"Name":"Test","CronLine":"5,15,25,35,45,55 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 2 $JOBID $TRIGGERTIME >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	SendJob(`{"Name":"Test","CronLine":"1,2,3,4,5,6,7,8,9 * * * * * *","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 3 $JOBID $TRIGGERTIME >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		for j := 0; j < 100000; j++ {
	// 			SendJob(`{"Name":"Test","CronLine":"1 * * * * * 2030","ExecutorType":"Shell","ExecutorInfo":"{\"Shell\":\"PowerShell\",\"Command\":\"echo 3 $JOBID $TRIGGERTIME >> test.txt\",\"Timeout\":\"10s\"}","RetryNum":1}`)
	// 		}
	// 	}()
	// }
	// for {
	// 	time.Sleep(10 * time.Second)
	// }
}
