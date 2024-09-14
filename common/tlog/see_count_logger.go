/*
	This is a sample class which show you how to implement the interface ICountLogger, and you can implement it by yourself
*/
package tlog

import (
	"context"
	"github.com/cihub/seelog"
	"log"
)

type SeeCountLogger struct {
	handler seelog.LoggerInterface
	counter ICounter
}

func InitCountLoggerFromSeelog(seeLogger seelog.LoggerInterface, counter ICounter) {
	log.Println("tg-flow CountLogger init start...")
	seeCountLogger := &SeeCountLogger{handler:seeLogger}
	Handler = NewCountLogger(seeCountLogger,counter)
	log.Println("tg-flow CountLogger init finished!!!")
}

func (d *SeeCountLogger) Debug(tag string, args ...interface{}) {
	d.handler.Debug(tag, args)
}

func (d *SeeCountLogger) Debugf(ctx context.Context, tag string, format string, args ...interface{}) {
	d.handler.Debugf(format, args)
}

func (d *SeeCountLogger) Info(tag string, args ...interface{}) {
	d.handler.Info(tag, args)
}

func (d *SeeCountLogger) Infof(ctx context.Context, tag string, format string, args ...interface{}) {
	d.handler.Infof(format, args)
}

func (d *SeeCountLogger) Error(tag string, args ...interface{}) {
	d.handler.Error(args)
}

func (d *SeeCountLogger) Errorf(ctx context.Context, tag string, format string, args ...interface{}) {
	d.handler.Errorf(format, args)
}

func (d *SeeCountLogger) Public(ctx context.Context, key string, pairs map[string]interface{}, isPLid bool) {
	// you can fullfill business log printer here
}
