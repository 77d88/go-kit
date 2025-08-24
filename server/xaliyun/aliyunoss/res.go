package aliyunoss

import (
	"time"

	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
)

const TableNameRes = "s_res"

const (
	ResTypeImage int32 = iota + 1 // 图片
	ResTypeVideo                  // 视频
	ResTypeAudio                  // 音频
	ResTypeFile                   // 文件
)

// Res mapped from table <s_res>
type Res struct {
	xpg.BaseModel           // 删除时间
	RefTime       time.Time `gorm:"comment:引用时间"`                   // 引用时间
	MimeType      int32     `gorm:"comment:资源类型"`                   // 资源类型
	Size          float64   `gorm:"comment:资源大小"`                   // 资源大小
	OptID         int64     `gorm:"comment:优化资源路径"`                 // 优化资源路径
	Cover         int64     `gorm:"comment:媒体资源类型-封面 如果是视频之类的"`     // 媒体资源类型-封面 如果是视频之类的
	Width         float64   `gorm:"comment:媒体资源类型-宽"`               // 媒体资源类型-宽
	Height        float64   `gorm:"comment:媒体资源类型-高"`               // 媒体资源类型-高
	AliEtag       string    `gorm:"comment:阿里eTag 文件一致性校验使用;index"` // 阿里eTag 文件一致性校验使用
	Path          string    `gorm:"not null;comment:路径"`            // 路径
	IsOptimize    bool      `gorm:"comment:是否优化;default:false"`     // 是否优化
}

// TableName Res's table name
func (*Res) TableName() string {
	return TableNameRes
}
