package local

import (
	"github.com/ericnts/log"
	"testing"
)

func TestAppend(t *testing.T) {
	data := make(map[string]map[string]string)
	m := make(map[string]string)
	m["success"] = "成功"
	data["ch"] = m
	Append(data)
	log.Infof(Get("ch", "success"))
	data2 := make(map[string]map[string]string)
	m2 := make(map[string]string)
	m2["success"] = "Success"
	data2["en"] = m2
	Append(data2)
	log.Infof(Get("ch", "success"))
	log.Infof(Get("en", "success"))

}
