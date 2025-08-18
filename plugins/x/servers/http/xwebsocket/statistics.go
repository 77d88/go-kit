package xwebsocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/77d88/go-kit/plugins/xlog"
)

type ConnectionStats struct {
	Timestamp     string          `json:"timestamp"`
	Total         int             `json:"totalConnections"`
	Groups        map[string]int  `json:"groups"`
	UserStats     UserConnections `json:"userStats"`
	TopConnection *TopConnector   `json:"topConnection,omitempty"`
	UserInfo      []UserConnInfo  `json:"userInfo,omitempty"`
}

type UserConnections struct {
	UniqueUsers int `json:"uniqueUsers"`
	MaxPerUser  int `json:"maxConnectionsPerUser"`
}

type TopConnector struct {
	UserID    int64 `json:"userId"`
	ConnCount int   `json:"connectionCount"`
}

type UserConnInfo struct {
	ClientID   int64  `json:"clientId"`
	UserID     int64  `json:"userId"`
	Group      string `json:"group"`
	Duration   string `json:"duration"`
	OnlineTime string `json:"onlineTime"`
}

// StartStatistics 统计
func (r *WsEngine) StartStatistics() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf(context.TODO(), "[websocket] 统计异常: %v", err)
			}
		}()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats := r.CollectStats()
			jsonData, _ := json.MarshalIndent(stats, "", "  ")
			xlog.Debugf(context.TODO(), "\nWebSocket 连接统计:\n%s", jsonData)

			// 生产环境可发送到监控系统
			//if !xapi.Server.Debug {
			//	sendToMonitoring(stats)
			//}
		}
	}()
}

// CollectStats 获取 engine 统计信息
func (r *WsEngine) CollectStats() ConnectionStats {
	r.statisticsMux.Lock()
	defer r.statisticsMux.Unlock()

	stats := ConnectionStats{
		Timestamp: time.Now().Format(time.RFC3339),
		Groups:    make(map[string]int),
	}

	// 基础统计
	stats.Total = len(r.clients)

	// 分组统计和用户映射
	userConnMap := make(map[int64]int)
	for _, v := range r.clients {
		stats.Groups[v.Group]++
		userConnMap[v.UserId]++
	}

	// 用户统计
	stats.UserStats.UniqueUsers = len(userConnMap)
	if stats.UserStats.UniqueUsers > 0 {
		maxConn := 0
		var topUser int64
		for uid, count := range userConnMap {
			if count > maxConn {
				maxConn = count
				topUser = uid
			}
		}
		stats.UserStats.MaxPerUser = maxConn
		stats.TopConnection = &TopConnector{
			UserID:    topUser,
			ConnCount: maxConn,
		}
	}

	debugInfos := make([]UserConnInfo, 0, len(r.clients))
	for _, v := range r.clients {
		duration := time.Since(v.ConnTime)
		debugInfos = append(debugInfos, UserConnInfo{
			ClientID:   v.ClientId,
			UserID:     v.UserId,
			Group:      v.Group,
			Duration:   duration.Round(time.Second).String(),
			OnlineTime: v.ConnTime.Format("2006-01-02 15:04:05"),
		})
	}
	stats.UserInfo = debugInfos
	return stats
}

//func sendToMonitoring(stats ConnectionStats) {
//	// 示例: 发送到Prometheus
//	metrics.Gauge("websocket_connections_total", float64(stats.Total))
//	for group, count := range stats.Groups {
//		metrics.Gauge("websocket_connections_by_group", float64(count), "group="+group)
//	}
//}
