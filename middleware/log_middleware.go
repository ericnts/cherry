package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/ericnts/cherry/current"
	"github.com/ericnts/cherry/results"
	"github.com/ericnts/config"
	"github.com/ericnts/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

type Record struct {
	Code            int           //HTTP状态吗
	StartTime       time.Time     //开始时间
	Latency         time.Duration //耗时
	CurrentUserID   string        //当前用户ID
	CurrentUserName string        //当前用户名称
	CurrentUserType string        //当前用户类型
	IP              string        //用户IP
	Agent           string        //用户代理
	URI             string        //请求URI
	Method          string        //请求方式
	Params          string        //请求参数
	Files           []FileMsg     //文件信息
	Response        string        //返回结果
	Error           string        //错误信息
}

// 文件信息
type FileMsg struct {
	Key  string
	Name string
	Size int64
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	// 获取response内容
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteHeader(code int) {
	switch code {
	case http.StatusFound:
		w.ResponseWriter.WriteHeader(code)
	default:
		// 统一返回200
		w.ResponseWriter.WriteHeader(http.StatusOK)
	}
}

func LogMiddleware(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	blw := bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	var params string
	var files []FileMsg
	if data, err := c.GetRawData(); err != nil {
		log.Error(err.Error())
	} else {
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		if form, err := c.MultipartForm(); form != nil && err == nil {
			if data, err := json.Marshal(form.Value); err != nil {
				params = string(data)
			}
			for key, headers := range form.File {
				fe := FileMsg{
					Key: key,
				}

				if headers != nil {
					fe.Name = headers[0].Filename
					fe.Size = headers[0].Size
				}

				files = append(files, fe)
			}
			postForm, _ := json.Marshal(c.Request.PostForm)
			params = string(postForm)
		} else {
			params = strings.ReplaceAll(strings.ReplaceAll(string(data), "\n", ""), " ", "")
		}
	}

	startTime := time.Now()

	c.Next() //执行
	//当前用户

	record := Record{
		StartTime:       startTime,
		IP:              c.ClientIP(),
		Agent:           c.Request.UserAgent(),
		URI:             c.Request.RequestURI,
		Method:          c.Request.Method,
		Params:          params,
		Files:           files,
		CurrentUserID:   current.UserID(c),
		CurrentUserName: current.UserName(c),
		CurrentUserType: current.UserType(c),
	}

	level := zapcore.DebugLevel
	if !c.IsWebsocket() && c.Writer.Header().Get("Content-Transfer-Encoding") != "binary" { //过滤ws和binary返回
		re := results.Result{}
		if err := json.Unmarshal(blw.bodyBuf.Bytes(), &re); err == nil {
			if re.Code >= 500 {
				level = zapcore.ErrorLevel
			} else if re.Code > 100 {
				level = zapcore.WarnLevel
			}
			if level != zapcore.DebugLevel || config.Options.LogResponse || record.Method != "GET" {
				record.Response = strings.Trim(blw.bodyBuf.String(), "\n")
			}
		}
	}
	record.Code = c.Writer.Status()                   //响应状态
	record.Latency = time.Now().Sub(record.StartTime) //运行时间
	errStrings := c.Errors.Errors()
	if errStrings != nil {
		record.Error = strings.Join(errStrings, ",")
	}

	switch level {
	case zapcore.DebugLevel:
		log.With("record", record).Debug(record.Method + " " + record.URI)
	case zapcore.InfoLevel:
		log.With("record", record).Info(record.Method + " " + record.URI)
	case zapcore.WarnLevel:
		log.With("record", record).Warn(record.Method + " " + record.URI)
	case zapcore.ErrorLevel:
		log.With("record", record).Error(record.Method + " " + record.URI)
	}

	c.Set("Record", record)

}
