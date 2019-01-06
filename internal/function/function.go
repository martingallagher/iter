package function

import (
	"reflect"
	"runtime"
	"strings"
)

// Name returns the name for the given function.
func Name(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	i := strings.LastIndexByte(name, '/')

	if i != -1 {
		name = name[i+1:]
	}

	return name
}
