package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

/* dkron的Shell执行器样例
{
	"executor": "shell",
	"executor_config": {
		"shell": "true",
		"command": "my_command",
		"env": "ENV_VAR=va1,ANOTHER_ENV_VAR=var2",
		"cwd": "/app",
		"timeout": "24h"
	}
}
*/

// Shell执行器
type ExecutorShell struct {
	Shell   string
	Command string
	Timeout string
	Env     string
}

func NewExecutorShell() *ExecutorShell {
	return &ExecutorShell{
		Timeout: "10s",
	}
}

func ParseExecutorShell(buf []byte) (*ExecutorShell, error) {
	res := &(ExecutorShell{})
	err := res.Parse(buf)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (e *ExecutorShell) Parse(buf []byte) error {
	err := json.Unmarshal(buf, e)
	if err != nil {
		return err
	}
	_, err = str2duration.ParseDuration(e.Timeout)
	return err
}

// 按照不同的Shell执行器获取前缀命令
// 目前只支持PowerShell和Bash
func getShellPrefix(shell string) []string {
	if strings.ToLower(shell) == "bash" {
		return []string{"bash", "-c"}
	}
	return []string{shell}
}

// 执行Shell命令，并且将幂等信息附于环境变量
func (e *ExecutorShell) Execute(jobId uint64, triggerTime time.Time) error {
	dur, err := str2duration.ParseDuration(e.Timeout)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	cmdLines := getShellPrefix(e.Shell)
	cmdLines = append(cmdLines, e.Command)
	cmd := exec.CommandContext(ctx, cmdLines[0], cmdLines[1:]...)

	newEnv := append(os.Environ(), fmt.Sprintf("JOBID=%d,TRIGGERTIME=%d", jobId, triggerTime.Unix()))
	if e.Env != "" {
		newEnv = append(newEnv, e.Env)
	}
	cmd.Env = newEnv

	err = cmd.Run()

	return err
}
