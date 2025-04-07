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
    em.mu.RLock()
    defer em.mu.RUnlock()

    // 构造正则表达式（将 * 转换为 .*）
    regexPattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
    reg, err := regexp.Compile(regexPattern)
    if err != nil {
        return
    }

    // 遍历所有已注册的事件名
    for eventName, handlers := range em.handlers {
        if reg.MatchString(eventName) {
            for _, h := range handlers {
                h(ctx)
            }
        }
    }
}
