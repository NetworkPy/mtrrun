package logger

type Logger interface {
	Info()
	Infof()
	Error()
	Errorf()
	Debug()
	Debugf()
	Fatal()
	Fatalf()
}
