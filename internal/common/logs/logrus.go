package logs

import (
	"github.com/stong1994/kit_golang/slog"
	"os"
	"strconv"
)

func Init() {
	isLocalEnv, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV"))
	slog.InitLog(true, 1, isLocalEnv)
	slog.SetLevel("debug")
}
