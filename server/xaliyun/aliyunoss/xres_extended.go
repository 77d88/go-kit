package aliyunoss

import (
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type LifeCycleOss struct {
	Client *oss.Client
	Config *Oss
}

func NewLifeCycleOss() *LifeCycleOss {
	return &LifeCycleOss{}
}

func (c *LifeCycleOss) Init(scanner xtype.Scanner) error {
	c.Client = InitWietScanner(scanner)
	return nil
}

func (c *LifeCycleOss) Dispose() error {
	return nil
}
