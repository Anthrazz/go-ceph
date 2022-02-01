// Package log allows to enable go-ceph logging and integrate it with the
// logging of the go-ceph consuming code.
package log

import (
	"sync/atomic"
	"unsafe"

	intLog "github.com/ceph/go-ceph/internal/log"
)

// go-ceph log levels
const (
	// Silence.
	NoneLvl = Level(intLog.NoneLvl)
	// Error events that might still allow the application to continue running.
	ErrorLvl = Level(intLog.ErrorLvl)
	// Potentially harmful situations.
	WarnLvl = Level(intLog.WarnLvl)
	// Informational messages that highlight the progress of the application at
	// coarse-grained level.
	InfoLvl = Level(intLog.InfoLvl)
	// Fine-grained informational events that are most useful to debug an
	// application.
	DebugLvl = Level(intLog.DebugLvl)
)

// Level of go-ceph logging.
type Level int32

// SetLevel sets the log level of the go-ceph logs.
//  PREVIEW
//
// The default log level is ErrorLvl.
func SetLevel(lvl Level) {
	atomic.StoreInt32(&intLog.Level, int32(lvl))
}

// Outputer must be implemented by the reveiver of the go-ceph logs.
type Outputer interface {
	Output(calldepth int, s string) error
}

// SetOutput sets the output destination for the logging.
//  PREVIEW
//
// The output can be set to any value, that implements the Outputer interface,
// which consist only of one output function. This interface is compatible to
// the standard log implementation, so that the easiest way to receive go-ceph
// logs, is to set it to the default logger:
//
//  import (
//    "log"
//    cephlog "github.com/ceph/go-ceph/common/log"
//  )
//
//  cephlog.SetOutput(log.Default())
//
// By default logs are disabled.
func SetOutput(o Outputer) {
	var p unsafe.Pointer
	if o != nil {
		outFunc := o.Output
		p = unsafe.Pointer(&outFunc)
	}
	atomic.StorePointer(&intLog.OutputPtr, p)
}

// assert that type internal.Output matches Outputer interface
func _() { _ = intLog.Output(Outputer(nil).Output) }
