package smallimgsave

import (
	"context"
	"fmt"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/external/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	xapi.InitTestConfig()
	xdb.Init()
	aliyunoss.Init(nil)

	client := aliyunoss.Client
	objectName := aliyunoss.Config.TempPrefix + "/" + xid.NextIdStr()
	// 生成PutObject的预签名URL
	result, err := client.Presign(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(aliyunoss.Config.OssBucket),
		Key:    oss.Ptr(objectName),
	},
		oss.PresignExpires(10*time.Minute),
	)
	if err != nil {
		t.Log("err", err)
		panic(err)
	}
	t.Log("result.url", result.URL)
	t.Log("result.SignedHeaders", result.SignedHeaders)
	t.Log("result.Expiration", result.Expiration)
	t.Log("result.Method", result.Method)

	uploadFile(result.URL, "G:\\Administrator\\Pictures\\f11f3a292df5e0fe8201d79f4e6034a85edf7216 (小).png")

}

func uploadFile(signedUrl, filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 创建一个新的HTTP客户端
	client := &http.Client{}

	// 创建一个PUT请求
	req, err := http.NewRequest("PUT", signedUrl, file)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("返回上传状态码: %d\n", resp.StatusCode)
	if resp.StatusCode == 200 {
		fmt.Println("使用网络库上传成功")
	}
	fmt.Println(string(body))

	return nil
}
