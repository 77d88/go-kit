package wxent

import (
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/plugins/xlog"
	"io"
	"net/http"
	"strings"
)

const (
	TestRobotKey = "4a26b3eb-5ee4-4c30-a429-6a649dc02d3c"
)

// EntRobotMsg 企业微信机器人
type EntRobotMsg struct {
	Key     string `json:"key"`     // 机器人ket
	Title   string `json:"title"`   // 标题 不超过128个字节，超过会自动截断
	Content string `json:"content"` // 消息内容
}

func NewEntRobotMsg(key, title string) *EntRobotMsg {
	return &EntRobotMsg{Key: key, Title: title}
}

// BuildArticles 构建图文消息
// url: 跳转链接
// picUrl: 图片链接
// description: 描述
func (m *EntRobotMsg) BuildArticles(url, picUrl, description string) *EntRobotMsg {
	msg := map[string]interface{}{
		"msgtype": "news",
		"news": map[string]interface{}{
			"articles": []map[string]interface{}{
				{
					"title":       m.Title,
					"description": description,
					"url":         url,
					"picurl":      picUrl,
				},
			},
		},
	}
	m.Content, _ = xparse.ToJSON(msg)
	return m
}
func (m *EntRobotMsg) BuildText(content string) *EntRobotMsg {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"Content": content,
		},
	}
	m.Content, _ = xparse.ToJSON(msg)
	return m
}
func (m *EntRobotMsg) BuildMarkdown(content string) *EntRobotMsg {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"Content": content,
		},
	}
	m.Content, _ = xparse.ToJSON(msg)
	return m
}
func (m *EntRobotMsg) Send() error {
	var url = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + m.Key

	resp, err := http.Post(url, "application/json", strings.NewReader(m.Content))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			xlog.Warnf(nil, "企微消息 关闭连接失败 %s", err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type Response struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	var respBody Response
	if err := json.Unmarshal(body, &respBody); err != nil {
		return err
	}
	if respBody.Errcode != 0 {
		return fmt.Errorf("failed to send message: %s", respBody.Errmsg)
	}

	return nil
}
