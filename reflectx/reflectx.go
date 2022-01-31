// Package reflectx provide helpers for reflect.
package reflectx

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// MethodsOf require pointer to interface (e.g.: new(app.Appl)) and
// returns all it methods.
func MethodsOf(v interface{}) []string {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		panic("require pointer to interface")
	}
	typ = typ.Elem()
	methods := make([]string, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		methods[i] = typ.Method(i).Name
	}
	return methods
}

// CallerPkgPath returns caller's package path for given stack depth.
func CallerPkgPath(skip int) string {
	return pkgPath(callerName(1 + skip))
}

// CallerPkg returns caller's package name (from path) for given stack depth.
func CallerPkg(skip int) string {
	return pkgName(callerName(1 + skip))
}

// CallerTypeMethodName returns caller's type.method name for given stack depth.
func CallerTypeMethodName(skip int) string {
	return stripTypeRef(typeMethodName(callerName(1 + skip)))
}

// CallerMethodName returns caller's method name for given stack depth.
func CallerMethodName(skip int) string {
	return methodName(callerName(1 + skip))
}

// CallerFuncName returns caller's func or method name for given stack depth.
func CallerFuncName(skip int) string {
	return funcName(callerName(1 + skip))
}

// Returns:
//   [example.com/path/]{dir|"main"}.{func|type.method}[."func"id[.id]...]
func callerName(skip int) string {
	pc, _, _, _ := runtime.Caller(1 + skip)
	return runtime.FuncForPC(pc).Name()
}

// Returns optional path plus dir or "main".
func pkgPath(name string) string {
	pos := strings.LastIndexByte(name, '/') + 1
	end := strings.IndexByte(name[pos:], '.')
	if end == -1 {
		panic(fmt.Sprintf("bad name: %s", name))
	} else {
		end += pos
	}
	return name[:end]
}

// Returns dir or "main".
func pkgName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	end := strings.IndexByte(name[start:], '.')
	if end == -1 {
		panic(fmt.Sprintf("bad name: %s", name))
	}
	return name[start : start+end]
}

// Returns type.method or func."func"id if it's not a method.
func typeMethodName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	pos := strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("bad name: %s", name))
	}
	start += pos + 1
	pos = strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("not a method name: %s", name))
	}
	end := strings.IndexByte(name[start+pos+1:], '.')
	if end == -1 {
		end = len(name)
	} else {
		end += start + pos + 1
	}
	return name[start:end]
}

func stripTypeRef(name string) string {
	if name[0] == '(' {
		pos := strings.IndexByte(name, ')')
		name = name[2:pos] + name[pos+1:]
	}
	return name
}

// Returns method or "func"id if it's not a method.
func methodName(name string) string {
	name = typeMethodName(name)
	pos := strings.IndexByte(name, '.')
	if pos == -1 {
		panic(fmt.Sprintf("not a method name: %s", name))
	}
	return name[pos+1:]
}

// Returns func or method if it's not a function.
func funcName(name string) string {
	start := strings.LastIndexByte(name, '/') + 1
	pos := strings.IndexByte(name[start:], '.')
	if pos == -1 {
		panic(fmt.Sprintf("bad name: %s", name))
	}
	start += pos + 1
	pos = strings.IndexByte(name[start:], '.')
	if pos == -1 {
		return name[start:]
	}
	return methodName(name)
}
