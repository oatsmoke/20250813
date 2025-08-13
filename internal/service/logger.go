package service

import (
	"context"
	"log"
)

type LoggerService struct {
	srvCtx context.Context
	msg    chan string
}

func NewLoggerService(ctx context.Context) *LoggerService {
	l := &LoggerService{
		srvCtx: ctx,
		msg:    make(chan string),
	}

	go l.start()

	return l
}

func (s *LoggerService) start() {
	for {
		select {
		case <-s.srvCtx.Done():
			return
		case msg := <-s.msg:
			log.Println(msg)
		}
	}
}

func (s *LoggerService) Print(msg string) {
	s.msg <- msg
}

func (s *LoggerService) Close() {
	close(s.msg)
}
