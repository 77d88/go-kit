package aliyunoss

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"hash"
	"io"
	"log"
	"time"
)

// {
// "region":"oss-cn-hangzhou",
// "endpoint":"https://oss-cn-hangzhou.aliyuncs.com",
// "accessKeyId":"xxx",
// "accessKeySecret":"xxx",
// "ossBucket":"oss-yzz-hangzhou",
// "domain":"https://xxx.xxx.cc",
// "savePrefix":"upload/yzz/res",
// "tempPrefix":"temp",
// "stsArn":"acs:ram::xxx:role/oss-manager",
// "stsEndpoint":"sts.cn-hangzhou.aliyuncs.com"
// }

type Oss struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	OssBucket       string `yaml:"ossBucket"`
	Domain          string `yaml:"domain"`
	StsEndpoint     string `yaml:"stsEndpoint"`
	StsArn          string `yaml:"stsArn"`
	Region          string `yaml:"region"`
	SavePrefix      string `yaml:"savePrefix"`
	TempPrefix      string `yaml:"tempPrefix"`
}

var (
	Client *oss.Client
	Config *Oss
)

func InitWith(x *x.Engine) *oss.Client {
	var config Oss
	x.Cfg.ScanKey("oss", &config)
	return Init(&config)
}

func Init(config *Oss) *oss.Client {
	Config = config

	region := Config.Region
	if region == "" {
		xlog.Errorf(nil, "oss config not found")
		return nil
	}

	if xstr.StartsWith(region, "oss-") { // 正常地域是不含oss-的 老版本的sdk 有
		Config.Region = region[4:]
		xlog.Warnf(nil, "oss region is v1 %s upgrade to %s", region, Config.Region)
	}

	xlog.Infof(nil, "load oss config: %s:%s:%s", Config.Region, Config.OssBucket, Config.Domain)
	c := oss.LoadDefaultConfig().
		WithRegion(Config.Region).
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(Config.AccessKeyId, Config.AccessKeySecret))
	Client = oss.NewClient(c)

	return Client
}

func GetOssPostSign(c context.Context) (map[string]interface{}, error) {
	// 构建Post Policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(1 * time.Hour) // 有效期一个小时
	region := Config.Region
	product := "oss"
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": Config.OssBucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request",
				Config.AccessKeyId, date, region, product)}, // 凭证
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			// 其他条件
			//[]any{"content-length-range", 1, 1024},
			// []any{"eq", "$success_action_status", "201"},
			// []any{"starts-with", "$key", "user/eric/"},
			// []any{"in", "$content-type", []string{"image/jpg", "image/png"}},
			// []any{"not-in", "$cache-control", []string{"no-cache"}},
		},
	}

	// 将Post Policy序列化为JSON字符串
	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalf("json.Marshal fail, err:%v", err)
		return nil, err
	}

	// 将Post Policy编码为Base64字符串
	stringToSign := base64.StdEncoding.EncodeToString(policy)

	// 生成签名密钥
	hmacHash := func() hash.Hash { return sha256.New() }
	signingKey := "aliyun_v4" + Config.AccessKeySecret
	h1 := hmac.New(hmacHash, []byte(signingKey))
	_, err = io.WriteString(h1, date)
	if err != nil {
		return nil, err
	}
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	_, err = io.WriteString(h2, region)
	if err != nil {
		return nil, err
	}
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	_, err = io.WriteString(h3, product)
	if err != nil {
		return nil, err
	}
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	_, err = io.WriteString(h4, "aliyun_v4_request")
	if err != nil {
		return nil, err
	}
	h4Key := h4.Sum(nil)

	// 计算Post签名
	h := hmac.New(hmacHash, h4Key)
	_, err = io.WriteString(h, stringToSign)
	if err != nil {
		return nil, err
	}
	signature := hex.EncodeToString(h.Sum(nil))

	return map[string]interface{}{
		"key":                     Config.TempPrefix + "/" + xid.NextIdStr(),
		"policy":                  stringToSign,
		"x-oss-signature-version": "OSS4-HMAC-SHA256",
		"x-oss-credential":        fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", Config.AccessKeyId, date, region, product),
		"x-oss-date":              utcTime.Format("20060102T150405Z"),
		"x-oss-signature":         signature,
		"url":                     fmt.Sprintf("https://%v.oss-%v.aliyuncs.com/", Config.OssBucket, region),
	}, nil
}
