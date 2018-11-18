package controllers

import (
	"sync"
	"time"

	"github.com/beego/ms304w-client/basis/conf"
	"github.com/satori/go.uuid"
)

var (
	AuthTimeout = conf.Int("auth_timeout")

	timeout = time.Minute * time.Duration(AuthTimeout)

	OAuth *Auth = NewAuth()
)

type Auth struct {
	lock *sync.RWMutex
	list map[string]time.Time
}

func NewAuth() *Auth {
	return &Auth{
		lock: new(sync.RWMutex),
		list: make(map[string]time.Time),
	}
}

func (m *Auth) Add() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	token := uuid.Must(uuid.NewV4()).String()
	m.list[token] = time.Now().Add(timeout)

	return token
}

func (m *Auth) Get(token string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if v, ok := m.list[token]; ok {
		if time.Now().After(v) {
			// del
			delete(m.list, token)

			return false
		}

		return true
	}

	return false
}

func (m *Auth) Set(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.list[token] = time.Now().Add(timeout)
}

func (m *Auth) Del(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.list, token)
}
