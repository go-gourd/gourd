package event

type HandlerEvent func(params any)

var event = make(map[string][]HandlerEvent)

// AddEvent 添加事件
func AddEvent(name string, callback HandlerEvent) {
	if _, ok := event[name]; ok {
		event[name] = append(event[name], callback)
	} else {
		event[name] = []HandlerEvent{
			callback,
		}
	}
}

// OnEvent 触发事件
func OnEvent(name string, params any) {
	if _, ok := event[name]; !ok {
		return
	}
	for i := range event[name] {
		event[name][i](params)
	}
	return
}
