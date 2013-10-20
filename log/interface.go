package log

type Interface interface {
	Emergency(format string, args ...interface{})
	Alert(format string, args ...interface{}) error
	Crit(format string, args ...interface{}) error
	Err(format string, args ...interface{}) error
	Warn(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}
