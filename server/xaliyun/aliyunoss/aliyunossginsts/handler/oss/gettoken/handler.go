package gettoken

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xcache"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

var StsClient *sts20150401.Client
var do sync.Once

func Init() {
	client, err := newStsClient()
	if err != nil {
		xlog.Errorf(nil, "sts client init fail %s", err)
		return
	}
	StsClient = client
}

func newStsClient() (*sts20150401.Client, error) {
	ossConfig := aliyunoss.Config
	config := &openapi.Config{AccessKeyId: &ossConfig.AccessKeyId, AccessKeySecret: &ossConfig.AccessKeySecret}
	// Endpoint 请参考 https://api.aliyun.com/product/Sts
	config.Endpoint = tea.String(ossConfig.StsEndpoint)
	_result := &sts20150401.Client{}
	_result, _err := sts20150401.NewClient(config)
	if _err != nil {
		return nil, _err
	}
	return _result, nil
}

// handler 获取token /oss/getToken
func handler(c *xhs.Ctx, r request) (interface{}, error) {
	var res response
	err := xcache.Once("oss:sts", &res, time.Minute*30, func() (interface{}, error) {
		assumeRoleRequest := &sts20150401.AssumeRoleRequest{
			DurationSeconds: xcore.V2p(int64(60 * 60)), // 设置STS令牌的有效期为1小时。
			RoleArn:         &aliyunoss.Config.StsArn,
			Policy: xcore.V2p(`{
                    "Version": "1",
                    "Statement": [{
                            "Action": ["oss:*"],
                            "Resource": ["acs:oss:*:*:*"],
                            "Effect": "Allow"
                        }]
                }`),
			RoleSessionName: xcore.V2p(fmt.Sprintf("oss@%s", uuid.New())),
		}
		// 设置运行时选项。
		runtime := &util.RuntimeOptions{}
		// 尝试执行获取STS令牌的操作，并处理可能的异常。
		res, tryErr := func() (r *sts20150401.AssumeRoleResponseBodyCredentials, _e error) {
			defer func() {
				if r := tea.Recover(recover()); r != nil {
					_e = r
				}
			}()
			// 调用StsClient的AssumeRoleWithOptions方法获取STS令牌。
			// 复制代码运行请自行打印 API 的返回值
			res, err := StsClient.AssumeRoleWithOptions(assumeRoleRequest, runtime)
			if err != nil {
				return nil, err
			}
			return res.Body.Credentials, nil
		}()

		// 如果尝试获取STS令牌时出现错误，进行错误处理。
		if tryErr != nil {

			var e = &tea.SDKError{}
			// 如果错误是tea.SDKError类型，则直接赋值。
			var _t *tea.SDKError
			if errors.As(tryErr, &_t) {
				e = _t
			} else {
				xlog.Errorf(c, "获取STS Token 失败 其他异常 %v", tryErr)
				return nil, tryErr
			}
			xlog.Errorf(c, "获取STS Token 失败 ali error %v", tryErr)
			fmt.Println(tea.StringValue(e.Message))
			// 尝试解析错误数据，获取推荐的解决措施。
			var data interface{}
			d := json.NewDecoder(strings.NewReader(tea.StringValue(e.Data)))
			err := d.Decode(&data)
			if err != nil {
				xlog.Errorf(c, "获取STS Token 失败 解析错误数据失败 %v", err)
				return nil, tryErr
			}
			if m, ok := data.(map[string]interface{}); ok {
				recommend, ok := m["Recommend"]
				if ok {
					xlog.Infof(c, "获取STS Token 失败 推荐的解决措施 %v", recommend)
				}
			}
			// 断言错误消息为字符串，如果失败则返回错误。
			_, err = util.AssertAsString(e.Message)
			if err != nil {
				xlog.Errorf(c, "获取STS Token 失败 断言错误消息为字符串失败 %v", err)
				return nil, err
			}
		}
		// 返回获取到的STS令牌信息。
		s := &response{
			Bucket:                            aliyunoss.Config.OssBucket,
			Region:                            aliyunoss.Config.Region,
			TempPrefix:                        aliyunoss.Config.TempPrefix,
			AssumeRoleResponseBodyCredentials: res,
		}
		return s, nil
	})
	if err != nil {
		xlog.Errorf(c, "获取token失败 %v", err)
		return nil, xerror.New("获取token失败")
	}
	return res, nil
}

// Run 获取token
func Run(c *xhs.Ctx) (interface{}, error) {
	do.Do(func() {
		Init()
	})
	var r request
	err := c.ShouldBind(&r)
	if err != nil {
		return nil, err
	}
	return handler(c, r)
}

type request struct {
}

type response struct {
	*sts20150401.AssumeRoleResponseBodyCredentials
	Bucket     string `json:"Bucket"`
	Region     string `json:"Region"`
	TempPrefix string `json:"TempPrefix"`
}
