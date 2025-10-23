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

func Listen(pattern string, callback Handler) {
	globalManager.Listen(pattern, callback)
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

// Listen 注册事件处理器，支持通配符模式
func (em *Manager) Listen(pattern string, callback Handler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[pattern] = append(em.handlers[pattern], callback)
}

// Trigger 触发事件
func (em *Manager) Trigger(name string, ctx context.Context) {
	// 先获取锁，复制匹配的处理器，然后释放锁
	var handlersToExecute []Handler
	em.mu.RLock()
	
	// 遍历所有注册的模式
	for pattern, handlers := range em.handlers {
		// 构造正则表达式（将 * 转换为 .*）
		regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
		reg, err := regexp.Compile(regexPattern)
		if err != nil {
			continue
		}
		
		if reg.MatchString(name) {
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