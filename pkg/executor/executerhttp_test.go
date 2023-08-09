package executor

import (
	"testing"
	"time"
)

const exampleExecutorHttpString string = `{` +
	`"Method":"GET",` +
	`"Url":"http://www.baidu.com",` +
	`"Headers":"{\"Hello\":\"World\"}",` +
	`"Body":"",` +
	`"Timeout":"10s",` +
	`"ExpectCode":200` +
	`}`

const badUrlExecutorHttpString string = `{` +
	`"Method":"GET",` +
	`"Url":"http://www.abadurl.com",` +
	`"Headers":"{}",` +
	`"Body":"",` +
	`"Timeout":"10s",` +
	`"ExpectCode":200` +
	`}`

const badJsonExecutorHttpString string = `{` +
	`"Method":"GET",` +
	`"Url":"http://www.abadurl.com",` +
	`"Headers":"{}",` +
	`"Body":"",` +
	`"Timeout":"10s",` +
	`"ExpectCode":200"` +
	`}`

const timeoutExecutorHttpString string = `{` +
	`"Method":"GET",` +
	`"Url":"http://www.baidu.com",` +
	`"Headers":"{\"Hello\":\"World\"}",` +
	`"Body":"",` +
	`"Timeout":"1ms",` +
	`"ExpectCode":200` +
	`}`

func TestHttpExecutor(t *testing.T) {
	res, err := ParseExecutorHttp([]byte(exampleExecutorHttpString))
	if err != nil {
		t.Errorf("good executor string, but fail to parse : %v", err)
	}
	err = res.Execute(1, time.Now())
	if err != nil {
		t.Errorf("good executor, but fail to execute : %v", err)
	}
	res, err = ParseExecutorHttp([]byte(badUrlExecutorHttpString))
	if err != nil {
		t.Errorf("good executor string, but fail to parse : %v", err)
	}
	err = res.Execute(1, time.Now())
	if err == nil {
		t.Errorf("bad executor, but succeed to execute")
	}
	_, err = ParseExecutorHttp([]byte(badJsonExecutorHttpString))
	if err == nil {
		t.Errorf("Bad executor string, but succeed to parse")
	}

	res, err = ParseExecutorHttp([]byte(timeoutExecutorHttpString))
	if err != nil {
		t.Errorf("good executor string, but fail to parse : %v", err)
	}
	err = res.Execute(1, time.Now())
	if err == nil {
		t.Errorf("should timeout, but succeed to execute : %v", err)
	}
}
