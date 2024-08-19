package views

var Router = map[string]func(){}

func RegisterRouterFunc(fnName string, fn func()) {
	Router[fnName] = fn
}
