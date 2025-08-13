package pgmq

import (
	"time"

	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

func init() {
	xdb.AddModels(&Queue{})
}

// Queue 队列消息结构
type Queue struct {
	xdb.BaseModel
	Message      string     `gorm:"column:message" json:"message,omitempty"`
	ReadTime     *time.Time `gorm:"column:read_time" json:"read_time,omitempty"`
	State        int16      `gorm:"column:state" json:"state"`
	Num          int16      `gorm:"column:num" json:"num"`
	Retry        int16      `gorm:"column:retry" json:"retry"`
	Type         int        `gorm:"column:type" json:"type"`
	ErrorInfo    string     `gorm:"column:error_info" json:"error_info,omitempty"`
	DeliveryTime time.Time  `gorm:"column:delivery_time" json:"delivery_time,omitempty"`
}

const TableNameQueue = "s_queue"

// TableName Res's table name
func (*Queue) TableName() string {
	return TableNameQueue
}
