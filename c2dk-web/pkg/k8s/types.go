package k8s

import "time"

var (
	count        int64 = 0
	errCount     int64 = 0
	timeoutCount int64
	connections  int           = 100
	timeout      time.Duration = time.Second * 5
	duration     time.Duration = time.Second * 10
)
