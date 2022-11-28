package factory

import (
	"io"
	"os"

	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers/jsonh"
	"github.com/zalgonoise/logx/handlers/texth"
)

func writerFromPath(path string) io.Writer {
	confFile, err := os.Open(path)
	if err != nil {
		newFile, err := os.Create(path)
		if err != nil {
			return os.Stderr
		}
		if err := newFile.Sync(); err != nil {
			return os.Stderr
		}
		return newFile
	}
	return confFile
}

func Logger(ltype, path string) logx.Logger {
	var (
		writer io.Writer = os.Stderr
		logger           = logx.New(texth.New(writer))
	)

	if path == "" {
		return logger
	}
	writer = writerFromPath(path)

	switch ltype {
	case "text":
		return logx.New(texth.New(writer))
	case "json":
		return logx.New(jsonh.New(writer))
	default:
		return logx.New(texth.New(writer))
	}
}
