package _interface

import (
	"fmt"
	"testing"
	"time"
)

type Interface interface {
	Start()
	Shutdown()
	Test(duration time.Duration)
}

type queue struct {
	Name string
	Data string
}

func (q *queue) Start() {
	fmt.Printf("Name: %s Start Data: %s\n", q.Name, q.Data)
}

func (q *queue) Shutdown() {
	fmt.Printf("Name: %s Stop Data: %s\n", q.Name, q.Data)
}

func (q *queue) Test(t time.Duration) {
	fmt.Println(t)
}

func TestInterFace(t *testing.T) {
	tests := []struct {
		Queue      *queue
		QueueStart func(Interface)
		QueueStop  func(Interface)
		Time       func(Interface, time.Duration)
	}{
		{
			Queue:      &queue{Name: "1", Data: "1"},
			QueueStart: Interface.Start,
			//QueueStop: func(i Interface) {},
			QueueStop: Interface.Shutdown,
			// Cannot use 'Interface.Test' (type func(Interface, time.Duration)) as the type func(Interface)
			//Time: Interface.Test,
			Time: Interface.Test,
		},
	}

	for _, v := range tests {
		v.QueueStart(v.Queue)
	}
}
