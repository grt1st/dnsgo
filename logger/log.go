package logger

import (
	"fmt"
	"log"
	"os"
)

const LOG_OUTPUT_BUFFER = 1024

const (
	LevelDebug = iota
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
)

type logMesg struct {
	Level int
	Mesg  string
}

type LoggerHandler interface {
	Setup(config map[string]interface{}) error
	Write(mesg *logMesg)
}

type Logger struct {
	level   int
	mesgs   chan *logMesg
	outputs map[string]LoggerHandler
}

func NewLogger() *Logger {
	logger := &Logger{
		mesgs:   make(chan *logMesg, LOG_OUTPUT_BUFFER),
		outputs: make(map[string]LoggerHandler),
	}
	go logger.Run()
	return logger
}

func (l *Logger) SetLogger(handlerType string, config map[string]interface{}) {
	var handler LoggerHandler
	switch handlerType {
	case "console":
		handler = NewConsoleHandler()
	case "file":
		handler = NewFileHandler()
	default:
		panic("Unknown log handler.")
	}

	handler.Setup(config)
	l.outputs[handlerType] = handler
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Run() {
	for {
		select {
		case mesg := <-l.mesgs:
			for _, handler := range l.outputs {
				handler.Write(mesg)
			}
		}
	}
}

func (l *Logger) writeMesg(mesg string, level int) {
	if l.level > level {
		return
	}

	lm := &logMesg{
		Level: level,
		Mesg:  mesg,
	}

	l.mesgs <- lm
}

func (l *Logger) Debug(format string, v ...interface{}) {
	mesg := fmt.Sprintf("[DEBUG] "+format, v...)
	l.writeMesg(mesg, LevelDebug)
}

func (l *Logger) Info(format string, v ...interface{}) {
	mesg := fmt.Sprintf("[INFO] "+format, v...)
	l.writeMesg(mesg, LevelInfo)
}

func (l *Logger) Notice(format string, v ...interface{}) {
	mesg := fmt.Sprintf("[NOTICE] "+format, v...)
	l.writeMesg(mesg, LevelNotice)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	mesg := fmt.Sprintf("[WARN] "+format, v...)
	l.writeMesg(mesg, LevelWarn)
}

func (l *Logger) Error(format string, v ...interface{}) {
	mesg := fmt.Sprintf("[ERROR] "+format, v...)
	l.writeMesg(mesg, LevelError)
}

type ConsoleHandler struct {
	level  int
	logger *log.Logger
}

func NewConsoleHandler() LoggerHandler {
	return new(ConsoleHandler)
}

func (h *ConsoleHandler) Setup(config map[string]interface{}) error {
	if _level, ok := config["level"]; ok {
		level := _level.(int)
		h.level = level
	}
	h.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return nil

}

func (h *ConsoleHandler) Write(lm *logMesg) {
	if h.level <= lm.Level {
		h.logger.Println(lm.Mesg)
	}
}

type FileHandler struct {
	level  int
	file   string
	logger *log.Logger
}

func NewFileHandler() LoggerHandler {
	return new(FileHandler)
}

func (h *FileHandler) Setup(config map[string]interface{}) error {
	if level, ok := config["level"]; ok {
		h.level = level.(int)
	}

	if file, ok := config["file"]; ok {
		h.file = file.(string)
		output, err := os.Create(h.file)
		if err != nil {
			return err
		}

		h.logger = log.New(output, "", log.Ldate|log.Ltime)
	}

	return nil
}

func (h *FileHandler) Write(lm *logMesg) {
	if h.logger == nil {
		return
	}

	if h.level <= lm.Level {
		h.logger.Println(lm.Mesg)
	}
}
