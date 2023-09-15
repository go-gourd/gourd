package event

import (
	"regexp"
	"strings"
)

type HandlerEvent func(params any)

// 存放注册的事件回调
var event = make(map[string][]HandlerEvent)

// Listen 监听事件
func Listen(name string, callback HandlerEvent) {
	if _, ok := event[name]; ok {
		event[name] = append(event[name], callback)
	} else {
		event[name] = []HandlerEvent{
			callback,
		}
	}
}

// Trigger 触发事件
func Trigger(name string, params any) {
	for eventName, handlers := range event {
		if matchEventName(eventName, name) {
			for _, handler := range handlers {
				handler(params)
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
