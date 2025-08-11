package xjob

import (
	"github.com/77d88/go-kit/plugins/xlog"
	"testing"
)

func TestName(t *testing.T) {
	manager, err := NewCronTaskManager()
	if err != nil {
		t.Fatal(err)
	}
	task := &CronTask{
		ID: "test",
		Job: func() error {
			xlog.Infof(nil, "test")
			return nil
		},
		Retry:   0,
		Spec:    "@every 1s",
		Timeout: 0,
	}
	_, err = manager.SubmitCronTask(task)
	if err != nil {
		t.Fatal(err)
	}
	select {}

}
