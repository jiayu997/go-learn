package log

import (
	"io"
	"log"
	"os"
)

var Info *log.Logger

func init(){
	file,err := os.OpenFile("console.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)

	if err != nil{
		log.Println("初始化日志文件失败.",err)
	}

	Info = log.New(io.MultiWriter(file,os.Stdout),"INFO:",log.Ldate|log.Ltime|log.Lshortfile)

}
