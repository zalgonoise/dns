package factory

import (
	"io"
	"os"

	"github.com/zalgonoise/logx"
	"github.com/zalgonoise/logx/handlers"
	"github.com/zalgonoise/logx/handlers/jsonh"
	"github.com/zalgonoise/logx/handlers/texth"
)

func writerFromPath(path string) io.Writer {
	confFile, err := os.Open(path)
	if err != nil {
		newFile, err := os.Create(path)
		if err != nil {
			return nil
		}
		if err := newFile.Sync(); err != nil {
			return nil
		}
		return newFile
	}
	return confFile
}

func Logger(ltype, path string) logx.Logger {
	var (
		defaultHandler = jsonh.New(os.Stderr)
		logger         = logx.New(defaultHandler)
	)

	if path == "" {
		return logger
	}
	fileWriter := writerFromPath(path)
	if fileWriter == nil {
		return logger
	}

	switch ltype {
	case "text":
		return logx.New(
			handlers.Multi(
				texth.New(fileWriter),
				defaultHandler,
			),
		)
	case "json":
		return logx.New(
			handlers.Multi(
				jsonh.New(fileWriter),
				defaultHandler,
			),
		)
	default:
		return logx.New(
			handlers.Multi(
				texth.New(fileWriter),
				defaultHandler,
			),
		)
	}
}
