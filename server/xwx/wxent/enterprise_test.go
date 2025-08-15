package wxent

import (
	"testing"

	"github.com/77d88/go-kit/plugins/xlog"
)

func TestEntRobotMsg_Send(t *testing.T) {

	err := NewEntRobotMsg(TestRobotKey, "测试").BuildArticles(
		"http://web.yuanzz.cc",
		"https://oss-img.b1.yuanzz.cc/cpup/2024-09-23/img/dacd217d-e83dfc23.jpg?x-oss-process=image/auto-orient,1/resize,m_lfit,w_200/quality,q_80",
		"测试").Send()
	if err != nil {
		xlog.Warnf(nil, "err %s", err)
		return
	}
}
