package event

import (
	"context"
	"regexp"
	"strings"
)

type Handler func(ctx context.Context)

// 存放注册的事件回调
var event = make(map[string][]Handler)

// Listen 监听事件
func Listen(name string, callback Handler) {
	if _, ok := event[name]; ok {
		event[name] = append(event[name], callback)
	} else {
		event[name] = []Handler{
			callback,
		}
	}
}

// Trigger 触发事件
func Trigger(name string, ctx context.Context) {
	for eventName, handlers := range event {
		if matchEventName(eventName, name) {
			for _, handler := range handlers {
				handler(ctx)
			}
		}
	}
	return
}

// 匹配事件名称
func matchEventName(eventName, pattern string) bool {
	if eventName == pattern {
		return true
	}
	escapedEventName := regexp.QuoteMeta(eventName)
	regexPattern := "^" + strings.Replace(pattern, "*", ".*", -1) + "$"

	matched, _ := regexp.MatchString(regexPattern, escapedEventName)
	return matched
}
