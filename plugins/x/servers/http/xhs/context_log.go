package xhs

func logFields(c *Ctx) map[string]interface{} {
	if c.test {
		return map[string]interface{}{
			"isLocal": true,
		}
	}

	var (
		method, reqUrl, clientIP string
	)

	if c.Request != nil {
		method = c.Request.Method
		reqUrl = c.Request.RequestURI
		clientIP = c.ClientIP()
	}

	return map[string]interface{}{
		"method": method,
		"uri":    reqUrl,
		"cip":    clientIP,
		"tid":    c.TraceId,
		"uid":    c.GetUserId(),
	}
}
