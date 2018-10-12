package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

var Term *logrus.Logger

// ContextHook ...
type ContextHook struct{}

// Levels ...
func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire ...
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(10); ok {
		funcName := runtime.FuncForPC(pc).Name()

		entry.Data["s"] = fmt.Sprintf("%16s:%4v:%-30v", path.Base(file), line, path.Base(funcName))
	}

	return nil
}

func setupTerm() {
	Term = logrus.New()
	Term.Out = os.Stdout

	Term.Formatter = new(TextFormatter)

	h := new(ContextHook)
	Term.AddHook(h)

	// Default starting log level
	Term.Level = logrus.WarnLevel
	Term.Level = logrus.DebugLevel
}

func init() {
	setupTerm()
}