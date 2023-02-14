package local

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ericnts/cherry/util"
	"github.com/ericnts/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"sync"
)

const (
	localeFileSuffix = ".json"

	LanguageKey = "Accept-Language"
	EN          = "en"
	ZH          = "zh"
	ZhHk        = "zh-HK"

	LocalePrivate = "private"
)

var (
	lock           sync.RWMutex
	localeMap      = make(map[string]map[string]string)
	localeFileDirs = [...]string{"./conf/locales", "../conf/locales", "../../conf/locales"}
)

type Local interface {
	LocalKey() string
}

type stringLocal string

func (s stringLocal) LocalKey() string {
	return string(s)
}

func init() {
	for _, localeFileDir := range localeFileDirs {
		files, err := ioutil.ReadDir(localeFileDir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if strings.Index(file.Name(), localeFileSuffix) < 0 {
				continue
			}
			fileBuffer, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", localeFileDir, file.Name()))
			if err != nil {
				panic(util.WrapErr(err, "语言包读取失败"))
			}
			data := make(map[string]map[string]string)
			if err := json.Unmarshal(fileBuffer, &data); err != nil {
				panic(util.WrapErr(err, "语言包解析失败"))
			}
			Append(data)
		}
		return
	}
	log.Warn("加载语言包失败")
}

func Append(data map[string]map[string]string) {
	lock.Lock()
	defer lock.Unlock()
	for ln, v := range data {
		if m, ok := localeMap[ln]; !ok {
			localeMap[ln] = v
		} else {
			for k, v := range v {
				m[k] = v
			}
		}
	}
}

func GetLanguage(ctx context.Context) string {
	var language string
	if c, ok := ctx.(*gin.Context); ok {
		language = c.GetHeader(LanguageKey)
	} else {
		value := ctx.Value(LanguageKey)
		if value != nil {
			language = value.(string)
		}
	}
	switch language {
	case EN, ZhHk:
		return language
	default:
		return ZH
	}
}

func Get(language, key string) string {
	lock.RLock()
	defer lock.RUnlock()
	var msg string
	if localeMap, ok := localeMap[language]; ok {
		msg = localeMap[key]
	}
	if msg == "" {
		msg = key
	}
	return msg
}

func Translate(ctx context.Context, key string) string {
	return Get(GetLanguage(ctx), key)
}
