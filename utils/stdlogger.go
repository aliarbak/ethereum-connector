package utils

import "github.com/rs/zerolog/log"

type StdLogger struct{}

func NewStdLogger() StdLogger {
	return StdLogger{}
}

func (s StdLogger) Print(v ...interface{}) {
	log.Trace().Msgf("%+v", v...)
}

func (s StdLogger) Printf(format string, v ...interface{}) {
	log.Trace().Msgf(format, v...)
}
func (s StdLogger) Println(v ...interface{}) {
	log.Trace().Msgf("%+v", v...)
}
