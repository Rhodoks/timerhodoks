package executor

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

/* dkron的Http执行器样例
{
	"executor": "http",
	"executor_config": {
		"method": "GET",
		"url": "http://example.com",
		"headers": "[]",
 		"body": "",
		"timeout": "30",
    	"expectCode": "200",
    	"expectBody": "",
    	"debug": "true"
	}
}
*/

// Http执行器
type ExecutorHttp struct {
	Method     string
	Url        string
	Headers    string
	Body       string
	Timeout    string
	ExpectCode int
}

func NewExecutorHttp() *ExecutorHttp {
	return &ExecutorHttp{
		Timeout: "10s",
	}
}

func ParseExecutorHttp(buf []byte) (*ExecutorHttp, error) {
	res := NewExecutorHttp()
	err := res.Parse(buf)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 从json解析执行器
func (e *ExecutorHttp) Parse(buf []byte) error {
	err := json.Unmarshal(buf, e)
	if err != nil {
		return err
	}
	_, err = str2duration.ParseDuration(e.Timeout)
	return err
}

// 执行Http请求，并且将幂等信息附于Header
func (e *ExecutorHttp) Execute(jobId uint64, triggerTime time.Time) error {
	dur, err := str2duration.ParseDuration(e.Timeout)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: dur}
	payload := strings.NewReader(e.Body)
	req, err := http.NewRequest(e.Method, e.Url, payload)

	if err != nil {
		return err
	}

	req.Header.Add("JOBID", strconv.Itoa(int(jobId)))
	req.Header.Add("TRIGGERTIME", strconv.Itoa(int(triggerTime.Unix())))

	headers := make(map[string]string)

	err = json.Unmarshal([]byte(e.Headers), &headers)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != e.ExpectCode {
		return errors.New("unexpected status code")
	}
	return nil
}
