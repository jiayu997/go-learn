#编译阶段
FROM registry.c2cloud.cn/library/golang:1.14.15-stretch

MAINTAINER tyf <yifan.tan@chinacreator.com>

#复制源码到GOPATH路径下
COPY ./src/. /go/src

#编译Linux可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -o main  /go/src/chinacreator.com/c2/license/licserver.go

#运行阶段,scratch是具有文件系统的空镜像
FROM scratch

# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=0 /go/main /

ENTRYPOINT ["/main"]