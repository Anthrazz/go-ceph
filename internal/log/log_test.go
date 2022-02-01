package log_test

import (
	"io"
	stdlog "log"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/ceph/go-ceph/common/log"
	intLog "github.com/ceph/go-ceph/internal/log"
	"github.com/stretchr/testify/assert"
)

func testLog() {
	intLog.Debugf("-%s-", "debug")
	intLog.Infof("-%s-", "info")
	intLog.Warnf("-%s-", "warn")
	intLog.Errorf("-%s-", "error")
}

var testOut = []string{
	"log_test.go:17: <go-ceph>[DBG]-debug-",
	"log_test.go:18: <go-ceph>[INF]-info-",
	"log_test.go:19: <go-ceph>[WRN]-warn-",
	"log_test.go:20: <go-ceph>[ERR]-error-",
	"",
}

func checkLines(t *testing.T, lines []string) {
	for i := range lines {
		assert.Equal(t, testOut[len(testOut)-len(lines)+i], lines[i])
	}
}

func captureAllOutput(f func()) string {
	oldout := os.Stdout
	olderr := os.Stderr
	oldlog := stdlog.Writer()
	defer func() {
		os.Stdout = oldout
		os.Stderr = olderr
		stdlog.SetOutput(oldlog)
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	stdlog.SetOutput(w)
	go func() {
		f()
		_ = w.Close()
	}()
	buf, _ := io.ReadAll(r)
	return string(buf)
}

func TestLogOffByDefault(t *testing.T) {
	out := captureAllOutput(func() { testLog() })
	assert.Empty(t, out)
}

func TestLogLevels(t *testing.T) {
	stdlog.Default()
	var out strings.Builder
	logger := stdlog.New(&out, "", stdlog.Lshortfile)
	log.SetOutput(logger)
	t.Run("WarnLvlByDefault", func(t *testing.T) {
		out.Reset()
		testLog()
		lines := strings.Split(out.String(), "\n")
		assert.Len(t, lines, 3)
		checkLines(t, lines)
	})
	t.Run("WarnLvl", func(t *testing.T) {
		out.Reset()
		log.SetLevel(log.WarnLvl)
		testLog()
		lines := strings.Split(out.String(), "\n")
		assert.Len(t, lines, 3)
		checkLines(t, lines)
	})
	t.Run("InfoLvl", func(t *testing.T) {
		out.Reset()
		log.SetLevel(log.InfoLvl)
		testLog()
		lines := strings.Split(out.String(), "\n")
		assert.Len(t, lines, 4)
		checkLines(t, lines)
	})
	t.Run("DebugLvl", func(t *testing.T) {
		out.Reset()
		log.SetLevel(log.DebugLvl)
		testLog()
		lines := strings.Split(out.String(), "\n")
		assert.Len(t, lines, 5)
		checkLines(t, lines)
	})
	log.SetOutput(nil)
}

func TestLoggerNotGC(t *testing.T) {
	var sb strings.Builder
	logger := stdlog.New(&sb, "none", 0)
	var dummy bool
	done := make(chan struct{})
	runtime.SetFinalizer(logger, func(interface{}) {
		t.Error("unreachable")
	})
	runtime.SetFinalizer(&dummy, func(interface{}) {
		close(done)
	})
	log.SetOutput(logger)
	runtime.GC()
	runtime.GC()
	<-done
}
