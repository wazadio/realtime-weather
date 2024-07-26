package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type messageObject struct {
	TimeStamp    string `json:"time_stamp"`
	Message      any    `json:"message"`
	FileName     string `json:"file_name"`
	FunctionName string `json:"function_name"`
	Line         int    `json:"line"`
	ID           any    `json:"ID"`
}

type LogEnvironment int

const (
	HTTP LogEnvironment = iota
	SCHEDULLER
)

type LogLevel string

const (
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

const (
	ERROR_RECOVER_INFROMATION string = "Unable to recover log information"
	reset                     string = "\033[0m"
	red                       string = "\033[31m"
	green                     string = "\033[32m"
	yellow                    string = "\033[33m"
	runtimeFilename           string = "runtime.log"
	cronFilename              string = "cron.log"
	logDir                    string = "logs/"
)

func InitLogger() {
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create directories: %s\n", err)
		return
	}

	runtimeLog, err := os.OpenFile(fmt.Sprintf("./%s%s", logDir, runtimeFilename), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("error creating runtime log : %s", err.Error())
	}
	defer runtimeLog.Close()

	cronLog, err := os.OpenFile(fmt.Sprintf("./%s%s", logDir, cronFilename), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("error creating runtime log : %s", err.Error())
	}
	defer cronLog.Close()
}

func Print(ctx context.Context, level LogLevel, message any) {
	var information strings.Builder
	pc, fileName, line, ok := runtime.Caller(1)
	if !ok {
		log.Println(ERROR_RECOVER_INFROMATION)
	}

	if level == ERROR {
		information.WriteString(fmt.Sprintf("%s%s%s: ", red, level, reset))
	} else if level == WARNING {
		information.WriteString(fmt.Sprintf("%s%s%s: ", yellow, level, reset))
	} else if level == INFO {
		information.WriteString(fmt.Sprintf("%s%s%s: ", green, level, reset))
	}

	fileNamePath := strings.Split(fileName, "/")
	fileName = fileNamePath[len(fileNamePath)-1]
	functionNamePath := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	functionName := functionNamePath[len(functionNamePath)-1]

	informationDetail, _ := json.Marshal(messageObject{
		TimeStamp:    time.Now().Format(time.DateTime),
		Message:      fmt.Sprintf("%+v", message),
		FunctionName: functionName,
		FileName:     fileName,
		Line:         line,
		ID:           ctx.Value("request_id"),
	})

	information.WriteString(string(informationDetail))
	fullMessage := information.String()
	if _, ok := ctx.Value("environment").(LogEnvironment); !ok || ctx.Value("environment").(LogEnvironment) == HTTP {
		f, err := os.OpenFile(fmt.Sprintf("%s%s", logDir, runtimeFilename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(fullMessage + "\n"); err != nil {
			log.Println(err)
		}
	} else if ctx.Value("environment") == SCHEDULLER {
		f, err := os.OpenFile(fmt.Sprintf("%s%s", logDir, cronFilename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(fullMessage + "\n"); err != nil {
			log.Println(err)
		}
	}

	log.Println(fullMessage)
}
