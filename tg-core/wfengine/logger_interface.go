package wfengine
/*
import (
	"bytes"
	"context"
	"log"
	"sync"
)

type Logger interface {
	Error(ctx context.Context, tag string, errMsg string)
}

var Log Logger
var logOnce sync.Once
func SetLogger(logger Logger){
	Log = logger
}

func ConfirmLoggerInit() {
	logOnce.Do(func(){
		if Log == nil {
			Log = &DefaultLogger{}
		}
	})
}

type DefaultLogger struct {}

func (i *DefaultLogger) Error(ctx context.Context, tag string, errMsg string){
	buf := bytes.Buffer{}
	buf.WriteString("[ERROR]||tag=")
	buf.WriteString(tag)
	buf.WriteString("||errMsg=")
	buf.WriteString(errMsg)
	log.Println(buf.String())
}*/