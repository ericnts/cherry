package jwt

import (
	"sync"
	"time"
)

const (
	// 超过100条记录后进行清理
	storeSize = 100
	// 记录保存一个月
	maxSaveTime int64 = 3600 * 24 * 30
)

type ISingleCheck interface {
	Logon(key string, t int64)
	Check(key string, t int64) bool
}

type memorySingleCheck struct {
	sync.RWMutex
	data map[string]int64
	max  int64
}

func (m *memorySingleCheck) Logon(key string, t int64) {
	m.Lock()
	m.data[key] = t
	now := time.Now().Unix()
	if len(m.data) > storeSize {
		for key, t := range m.data {
			if now-t > maxSaveTime {
				delete(m.data, key)
			}
		}
	}
	m.Unlock()
}

func (m *memorySingleCheck) Check(key string, t int64) bool {
	m.RLock()
	value, ok := m.data[key]
	m.RUnlock()
	if ok {
		if t < value {
			return false
		}
	} else {
		m.Logon(key, t)
	}
	return true
}
