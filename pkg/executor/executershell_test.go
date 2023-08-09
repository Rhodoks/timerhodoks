package executor

import (
	"testing"
	"time"
)

const exampleShellExecutorString string = `{"Shell":"PowerShell","Command":"dir","Timeout":"2s"}`
const timeoutShellExecutorString string = `{"Shell":"PowerShell","Command":"sleep 10s","Timeout":"2s"}`
const wrongJsonShellExecutorString string = `{"Shell":"PowerShell","Command":"dir"","Timeout":"2s"}`
const wrongDurationShellExecutorString string = `{"Shell":"PowerShell","Command":"dir","Timeout":"2s1"}`

func TestShellExecutor(t *testing.T) {
	res, err := ParseExecutorShell([]byte(exampleShellExecutorString))
	if err != nil {
		t.Errorf("exampleShellExecutorString: Right json, but fail to parse")
	}
	err = res.Execute(1, time.Now())
	if err != nil {
		t.Errorf("exampleShellExecutorString: Fail to run (%v)", err)
	}

	res, err = ParseExecutorShell([]byte(timeoutShellExecutorString))
	if err != nil {
		t.Errorf("timeoutShellExecutorString: Right json, but fail to parse.")
	}
	err = res.Execute(1, time.Now())
	if err == nil {
		t.Errorf("timeoutShellExecutorString: Should timeout, but does not.")
	}

	_, err = ParseExecutorShell([]byte(wrongJsonShellExecutorString))
	if err == nil {
		t.Errorf("wrongDurationShellExecutorString: Wrong json, but succeed to parse.")
	}

	_, err = ParseExecutorShell([]byte(wrongDurationShellExecutorString))
	if err == nil {
		t.Errorf("wrongDurationShellExecutorString: Wrong duration, but succeed to parse.")
	}
}
