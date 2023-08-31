package event

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
	if _, ok := event[name]; !ok {
		return
	}
	for i := range event[name] {
		event[name][i](params)
	}
	return
}
