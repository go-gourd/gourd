package event

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

type Handler func(ctx context.Context)

// 全局单例
var globalManager = NewEventManager()

func Listen(name string, callback Handler) {
	globalManager.Listen(name, callback)
}

func Trigger(name string, ctx context.Context) {
	globalManager.Trigger(name, ctx)
}

type Manager struct {
	handlers map[string][]Handler // 存储所有事件处理器（key为注册时的原始名称）
	mu       sync.RWMutex
}

func NewEventManager() *Manager {
	return &Manager{
		handlers: make(map[string][]Handler),
	}
}

// Listen 注册事件处理器
func (em *Manager) Listen(name string, callback Handler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[name] = append(em.handlers[name], callback)
}

// Trigger 触发事件
func (em *Manager) Trigger(pattern string, ctx context.Context) {
	// 构造正则表达式（将 * 转换为 .*）
	regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
	reg, err := regexp.Compile(regexPattern)
	if err != nil {
		return
	}

	// 先获取锁，复制匹配的处理器，然后释放锁
	var handlersToExecute []Handler
	em.mu.RLock()
	for eventName, handlers := range em.handlers {
		if reg.MatchString(eventName) {
			// 复制处理器到临时切片
			handlersToExecute = append(handlersToExecute, handlers...)
		}
	}
	em.mu.RUnlock()

	// 在不持有锁的情况下执行处理器
	for _, h := range handlersToExecute {
		h(ctx)
	}
}
