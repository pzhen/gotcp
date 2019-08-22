package Gotcp

import (
	"fmt"
	"gotcp/Conf"
	"os"
	"strings"
	"time"
)

const (
	LevelInfo = "\x1b[94m[INFO-%s]\x1b[0m"
	LevelTrac = "\x1b[93m[TRAC-%s]\x1b[0m"
	LevelErro = "\x1b[91m[ERRO-%s]\x1b[0m"
	LevelWarn = "\x1b[95m[WARN-%s]\x1b[0m"
	LevelSucc = "\x1b[92m[SUCC-%s]\x1b[0m"
)

func debugPrint(format string, values ...interface{}) {
	print(LevelInfo, format, values ...)
}

func debugPrintTrace(format string, values ...interface{}) {
	print(LevelTrac, format, values ...)
}

func debugPrintWarn(format string, values ...interface{}) {
	print(LevelWarn, format, values ...)
}

func debugPrintSucc(format string, values ...interface{}) {
	print(LevelSucc, format, values ...)
}

func debugPrintError(format string, values ...interface{}) {
	print(LevelErro, format, values ...)
}

func print(level string, format string, values ...interface{})  {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	prefix := time.Now().Format("2006/01/02 15:04:05") + " " + fmt.Sprintf(level, Conf.SrvConf.Env) + " "
	fmt.Fprintf(os.Stderr, prefix + format, values...)
}
