// Package log is the internal package for go-ceph logging. This package is only
// used from go-ceph code, not from consumers of go-ceph. Go code uses the log
// functions like Errorf() etc. to log information, and the consumer of go-ceph
// decides with SetOutput() and SetLevel() from the exported log package, if and
// what of these logs they want to receive.
package log

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// Log levels
const (
	NoneLvl = iota
	ErrorLvl
	WarnLvl
	InfoLvl
	DebugLvl
)

const (
	gocephPrefix = "<go-ceph>"
	errorPrefix  = "[ERR]"
	warnPrefix   = "[WRN]"
	infoPrefix   = "[INF]"
	debugPrefix  = "[DBG]"
)

const defaultLvl = WarnLvl

// These variables are set by the common log package.
var (
	OutputPtr unsafe.Pointer // pointer to type Output
	Level     int32          = defaultLvl
)

// Output is the signature of the Output function.
type Output func(calldepth int, s string) error

// Errorf logs events that might still allow the application to continue
// running. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	o := getOut()
	if o == nil {
		return
	}
	logOut(*o, ErrorLvl, format, v)
}

// Warnf logs potentially harmful situations. Arguments are handled in the
// manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	o := getOut()
	if o == nil {
		return
	}
	logOut(*o, WarnLvl, format, v)
}

// Infof logs informational messages that highlight the progress of the
// application at coarse-grained level. Arguments are handled in the manner of
// fmt.Printf.
func Infof(format string, v ...interface{}) {
	o := getOut()
	if o == nil {
		return
	}
	logOut(*o, InfoLvl, format, v)
}

// Debugf logs fine-grained informational events that are most useful to debug
// an application. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	o := getOut()
	if o == nil {
		return
	}
	logOut(*o, DebugLvl, format, v)
}

func logOut(out Output, typ int32, format string, v []interface{}) {
	l := getLvl()
	if l < typ {
		return
	}
	_ = out(3, fmt.Sprintf(prefix(typ)+format, v...))
}

func getOut() *Output {
	return (*Output)(atomic.LoadPointer(&OutputPtr))
}

func getLvl() int32 {
	return atomic.LoadInt32(&Level)
}

func prefix(lvl int32) string {
	var prefix string
	switch lvl {
	case ErrorLvl:
		prefix = errorPrefix
	case WarnLvl:
		prefix = warnPrefix
	case InfoLvl:
		prefix = infoPrefix
	case DebugLvl:
		prefix = debugPrefix
	}
	return gocephPrefix + prefix
}
