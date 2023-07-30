package k8s

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DnsCheck(host, dnsserver string) string {
	pool := make(chan interface{}, connections)
	exit := make(chan bool)
	var (
		min int64 = 0
		max int64 = 0
		sum int64 = 0
	)

	go func() {
		time.Sleep(duration)
		exit <- true
	}()

endD:
	for {
		select {
		case pool <- nil:
			go func() {
				defer func() {
					<-pool
				}()
				now := time.Now()
				resolver := &net.Resolver{
					PreferGo: true,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						d := net.Dialer{
							Timeout: timeout * time.Second,
						}
						return d.DialContext(ctx, "udp", dnsserver+":53")
					},
				}
				_, err := resolver.LookupIPAddr(context.Background(), host)
				//use := time.Since(now).Milliseconds() / int64(time.Millisecond)
				use := time.Since(now).Milliseconds()
				if min == 0 || use < min {
					min = use
				}
				if use > max {
					max = use
				}
				sum += use
				if use >= time.Duration.Milliseconds(timeout) {
					//timeoutCount++
					atomic.AddInt64(&timeoutCount, 1)
				}
				atomic.AddInt64(&count, 1)
				if err != nil {
					//fmt.Println(err.Error())
					atomic.AddInt64(&errCount, 1)
				}
			}()
		case <-exit:
			break endD
		}
	}

	//fmt.Printf("request count：%d\nerror count：%d\n", count, errCount)
	//fmt.Printf("request time：min(%dms) max(%dms) avg(%dms) timeout(%dn)\n", min, max, sum/count, timeoutCount)
	return (fmt.Sprintf("请求总数：%d 错误数：%d 请求耗时：min(%dms)--max(%dms)--avg(%dms)--timeout(%dn)\n", count, errCount, min, max, sum/count, timeoutCount))
}

// 返回service IP
func GetDNSER(namespace, name string) (string, error) {
	client, err := initClientSet()
	if err != nil {
		return "K8S连接失败", err
	}
	service, err := client.CoreV1().Services(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}
