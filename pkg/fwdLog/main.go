package fwdLog

import (
	llog "log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type fwdLogger struct {
	logger *zerolog.Logger
}

func (l *fwdLogger) Write(p []byte) (n int, err error) {
	log.Error().
		Str("error", string(p)).
		Msg("http server error")

	return len(p), nil
}

func (l *fwdLogger) log() *llog.Logger {
	return llog.New(&fwdLogger{}, "", 0)
}

func Logger() *llog.Logger {
	var log fwdLogger
	return log.log()
}
