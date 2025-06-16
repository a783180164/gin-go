package logger

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// Log 是全局可用的 logrus Logger 实例
var Log *logrus.Logger

// InitDateLogger 初始化全局 logrus 日志实例，并开启文件切割。
// pattern  ：日志文件命名格式，如 "/var/log/myapp.%Y%m%d.log"
// linkName  ：为方便获取最新日志，可创建一个指向最新日志的符号链接，如 "/var/log/myapp.log"
// retainDays：日志保留天数，超过会被自动删除
// level     ：日志级别，如 logrus.InfoLevel、logrus.DebugLevel 等
func InitDateLogger(pattern, linkName string, retainDays int, level logrus.Level) (*logrus.Logger, error) {
	writer, err := rotatelogs.New(
		pattern,
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithMaxAge(time.Duration(retainDays)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetOutput(writer)
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})
	Log = logger
	return logger, nil
}

// bodyLogWriter 用于捕获响应 body 和大小
type bodyLogWriter struct {
	gin.ResponseWriter
	body         *bytes.Buffer
	responseSize int
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// 在写入客户端之前，也将响应内容写入到 body 缓存
	w.body.Write(b)
	n, err := w.ResponseWriter.Write(b)
	w.responseSize += n
	return n, err
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	// 如果调用 WriteString，也要写到 buffer 里
	w.body.WriteString(s)
	n, err := w.ResponseWriter.WriteString(s)
	w.responseSize += n
	return n, err
}

// GinLogger 返回一个 Gin 中间件，用于记录请求和响应日志。
//
//	enableRequestLog  ：是否记录“请求日志”（请求进入时即记录请求方法、路径、客户端 IP、User-Agent、请求体长度等）
//	enableResponseLog ：是否记录“响应日志”（处理完成后记录响应状态码、响应体（最多 4KB）、耗时等）
//
// 如果想分别在不同场景开启/关闭，只需传入不同的布尔值即可。
func GinLogger(logger *logrus.Logger, enableRequestLog bool, enableResponseLog bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 捕获请求数据
		method := c.Request.Method
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 如果需要记录请求日志，则读取请求 Body 大小（注意：这里没有读取完整 Body；若需要记录完整请求 Body，需要额外处理）
		// 这里只记录 ContentLength，表示请求体大小
		requestSize := c.Request.ContentLength

		// 用自定义的 bodyLogWriter 替换 c.Writer，后续写入会被缓存到 bodyLogWriter.body
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			responseSize:   0,
		}
		c.Writer = blw

		// 在处理请求前，如果开启了请求日志，先输出一条“请求”日志
		if enableRequestLog {
			entry := logger.WithFields(logrus.Fields{
				"method":       method,
				"path":         path,
				"query":        rawQuery,
				"ip":           clientIP,
				"user_agent":   userAgent,
				"request_size": requestSize,
			})
			entry.Info("request started")
		}

		// 继续后续处理
		c.Next()

		// 记录处理耗时、状态码等
		latency := time.Since(start)
		statusCode := blw.Status()

		// 如果开启了响应日志，则从 blw.body 中获取响应内容
		var respBodySnippet string
		const maxBodyLogSize = 4 * 1024 // 最多 4KB
		if enableResponseLog {
			fullBody := blw.body.Bytes()
			if len(fullBody) > maxBodyLogSize {
				respBodySnippet = string(fullBody[:maxBodyLogSize]) + "...(truncated)"
			} else {
				respBodySnippet = string(fullBody)
			}
		}

		// 最终构造一个日志条目，包含通用字段 + 可选的响应内容
		fields := logrus.Fields{
			"status":      statusCode,
			"latency":     latency,
			"response_sz": blw.responseSize,
		}
		if enableResponseLog {
			fields["response_body"] = respBodySnippet
		}

		entry := logger.WithFields(fields)
		entry.Info("request completed")
	}
}
