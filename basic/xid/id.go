package xid

import (
	"fmt"
	"sync"
	"time"
)

// WorkerOsMethod 支持两种生成方式 1. 固定WorkerId 2. 获取当前操作系统的workerId
const WorkerOsMethod = "ID_WORKER_ID"
const WorkerOsParam = "ID_WORKER_ID"

var singletonMutex sync.Mutex
var idGenerator *DefaultIdGenerator

// SetIdGenerator .
func SetIdGenerator(options *IdGeneratorOptions) {
	singletonMutex.Lock()
	idGenerator = NewDefaultIdGenerator(options)
	singletonMutex.Unlock()
}

// NextId .
func NextId() int64 {
	if idGenerator == nil {
		SetIdGenerator(NewIdGeneratorOptions(1))
	}

	return idGenerator.NewLong()
}

func ExtractTime(id int64) time.Time {
	return idGenerator.ExtractTime(id)
}

func NextIdStr() string {
	return fmt.Sprintf("%d", NextId())
}
