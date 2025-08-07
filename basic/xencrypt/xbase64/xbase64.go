package xbase64

import "encoding/base64"

func Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func EncodeURL(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func DecodeURL(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

// RawURLEncode 去除填充字符版本 需要url传输优先使用这个
func RawURLEncode(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

func RawURLDecode(src string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(src)
}
