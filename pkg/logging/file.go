package logging

import (
	"fmt"
	"os"
	"strconv"
	"time"

	c "whisper/pkg/configuration_center/client"
)

// getLogFilePath get the log file save path

// Logging log 配置
type Logging struct {
	LogSavePath     string `json:"logSavePath"`
	LogSaveName     string `json:"logSaveName"`
	LogFileExt      string `json:"logFileExt"`
	TimeFormat      int    `json:"timeFormat"`
	RuntimeRootPath string `json:"runtimeRootPath"`
}

func getLogFilePathAndName() (string, string) {
	var result = c.C()
	var cfg = &Logging{}
	var _ = result.App("logging", cfg)
	gopath := os.Getenv("GOPATH")
	fmt.Println("cfg.Logging.RuntimeRootPath--->", cfg.RuntimeRootPath)
	return fmt.Sprintf("%s%s%s%s", gopath, "/src/whisper/", cfg.RuntimeRootPath, cfg.LogSavePath),
		fmt.Sprintf("%s%s.%s",
			cfg.LogSaveName,
			time.Now().Format(strconv.Itoa(cfg.TimeFormat)),
			cfg.LogFileExt)
}
