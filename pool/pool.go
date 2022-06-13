package pool

import (
	"github.com/ZinkLu/TGRobot/handlers/common"
)

type appHandlersMap map[string]common.AppHandlerInterface

var POOL = make(appHandlersMap, 0)

func Add(handlers ...common.AppHandlerInterface) {
	for _, handler := range handlers {
		name := handler.Name()
		_, ok := POOL[name]
		if ok {
			panic("name:" + name + "exists!")
		}
		POOL[name] = handler
	}
}

// name is unnecessary maybe..
func GetAppName[T common.AppHandlerInterface](name string) T {
	h, ok := POOL[name]
	if !ok {
		panic("no handler named " + name + "has register to bot")
	}
	return h.(T) // panic if h nil ?
}
