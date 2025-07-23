package xdbutil

import (
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"strconv"
	"strings"
)

type ConnectionInfo struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

// ToDns 链接转DNS 如 host=pgsql port=5432 user=postgres password=123456 dbname=hospital sslmode=disable
func (c ConnectionInfo) ToDns(aps ...string) string {
	return strings.Join([]string{
		"host=" + c.Host,
		"port=" + strconv.Itoa(c.Port),
		"user=" + c.Username,
		"password=" + c.Password,
		"dbname=" + c.Database,
		strings.Join(aps, " "),
	}, " ")
}

func (c ConnectionInfo) ToIntDatabase() int {
	atoi, err := strconv.Atoi(c.Database)
	if err != nil {
		fmt.Println()
		xlog.Fatalf(nil, "redis config error:%s", err)
		return 0
	}
	return atoi
}

// ParseConnection 链接解析
// 如 pgsql://postgres:123456@127.0.0.1:5432/hospital
// 如 redis://:123456@127.0.0.1:6379/0
func ParseConnection(connStr string) (*ConnectionInfo, error) {
	info := &ConnectionInfo{Port: -1}

	// 分离协议部分
	if protoSep := strings.Index(connStr, "://"); protoSep > 0 {
		info.Protocol = connStr[:protoSep]
		connStr = connStr[protoSep+3:]
	}

	// 提取认证信息（增强部分）
	if atSep := strings.Index(connStr, "@"); atSep > 0 {
		authPart := connStr[:atSep]
		connStr = connStr[atSep+1:]

		// 处理只有密码的情况（如redis://:password@host）
		if strings.HasPrefix(authPart, ":") {
			info.Password = authPart[1:]
		} else if colonSep := strings.Index(authPart, ":"); colonSep > 0 {
			info.Username = authPart[:colonSep]
			info.Password = authPart[colonSep+1:]
		} else {
			info.Username = authPart
		}
	}

	// 处理主机和端口
	if slashSep := strings.Index(connStr, "/"); slashSep > 0 {
		hostPort := connStr[:slashSep]
		info.Database = connStr[slashSep+1:]

		if colonSep := strings.Index(hostPort, ":"); colonSep > 0 {
			info.Host = hostPort[:colonSep]
			if port, err := strconv.Atoi(hostPort[colonSep+1:]); err == nil {
				info.Port = port
			}
		} else {
			info.Host = hostPort
		}
	} else {
		if colonSep := strings.Index(connStr, ":"); colonSep > 0 {
			info.Host = connStr[:colonSep]
			if port, err := strconv.Atoi(connStr[colonSep+1:]); err == nil {
				info.Port = port
			}
		} else {
			info.Host = connStr
		}
	}

	// 设置默认端口
	if info.Port == -1 {
		switch info.Protocol {
		case "redis":
			info.Port = 6379
		case "postgres", "postgresql":
			info.Port = 5432
		case "mysql":
			info.Port = 3306
		case "mssql", "sqlserver":
			info.Port = 1433
		case "oracle":
			info.Port = 1521
		case "mongodb":
			info.Port = 27017
		}
	}

	return info, nil
}
