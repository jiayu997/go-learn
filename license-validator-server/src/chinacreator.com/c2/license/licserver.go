package main

import (
	"chinacreator.com/c2/license/log"
	rsa "chinacreator.com/c2/license/utils"
	"chinacreator.com/c2/license/validator"
	"chinacreator.com/c2/license/validator/cpu"
	"chinacreator.com/c2/license/validator/disk"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type License struct {
	Id				 string
	Mac              []string
	Cpuid            string
	DiskSerialNumber string
	RequestTime      string
}

type Accept struct {
	AcceptSN []string
	Id 	string
}

var port string

func init(){
	//通过环境变量获取port配置
	if port=="" {
		licport := os.Getenv("LIC_PORT")
		if licport != "" {
			port = licport
		}
	}

}

var privateKey []byte

func init(){
	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDD9xgCkCmtZ4ZM1KTYgerm+H3yNBltpf9EYYemxtlLCoG1EawK
8SpRiAVBa1CudPppJmE2wDi1TfcemwkYfKFe2dHoySx7+3J4Mbd1uyyTWyt9ebBa
snC2Mbao3mtiTzpvunzIPbfpSEMniNtoYlM1knuWC9fHpEMW+80Oq2pbawIDAQAB
AoGANhnIciS0rN/QzvNB01gCruNZef1yK7hRQeKfHab2JGZxKrkHQzoTUdD4ingD
HTbETzU+T2w/+6XbnIJ2v2Dg96AEcZxs3hZ3V+QNlCLXJOk2s3lZp9646X4oAJcF
UDrifUImgLuQZV3Qgi9fJEgbi1MIhZgq3w6ZBCokdLIfBKkCQQDensuBRTWmKL82
Gq1cA/5bNvPheBMvgiZjSuHyTlpdsbgrpNhkZoWKpHbcQNguPkR6badaQucLAZN0
Jg5EJ/n3AkEA4VknCA5KboJ/eakIfVWqJc9CS+XND0QTcuIjVHNLbpaCqNa8YzpB
vMXP/OjXU9lajaZhWS93qYMgab8ZjWUtLQJBAMFP6O1m+PBBX9EOl01Y1m3EqUA3
sYlGnikIpG1xZn0HzyJu8c01TW8X43LdCBwXzAT35SO3BsQC6VUpmqfKgv8CQBdH
UGLipwm3bVeyAHCCEuuI935DpOU40RGDDsdAicBIyAKM/DT75aKMhKnJm8TLpTEQ
yOmfn6rhIs4JsagLlZkCQC3GGd9axIs/FkvQLrumbOEnYBfSVZDZGzce4WjwL2E5
Mbh2t4Bl0jZZ0PTz0JH1UXuHwyIeVj/NU+cwoOjc4gM=
-----END RSA PRIVATE KEY-----`)

}

/*func readBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}*/

func handler(writer http.ResponseWriter, request *http.Request) {
	//LICENSEID用于提供追踪验证请求来源License
	licsn := request.Header.Get("LIC-ID")
	log.Info.Println(licsn+"接收到请求")
	if request.Method == "POST" {
		//读取lic验证请求内容
		body, err := ioutil.ReadAll(request.Body)
		if len(body)==0 {
			log.Info.Println(licsn+"请求参数无效")
			http.Error(writer, licsn+"请求参数无效", 500)
			return
		}
		if err != nil {
			log.Info.Println(licsn+"请求读取请求内容出错",err)
			http.Error(writer, licsn+"请求读取请求内容出错", 500)
			return
		}
		d, deErr := rsa.RsaDecrypt(privateKey,body)

		if deErr != nil{
			log.Info.Println(licsn+"请求解码报错,请联系管理员",deErr)
			http.Error(writer, licsn+"请求解码报错,请联系管理员", 500)
			return
		}
		lic := new(License)
		_ = json.Unmarshal(d, &lic)

		//验证请求时间差,默认为60秒内有效
		now := time.Now()
		lictime,_ := strconv.ParseInt(lic.RequestTime, 10, 64)
		diff := now.Sub(time.Unix(0, lictime * int64(time.Millisecond)))

		if diff.Seconds() > 600 || diff.Seconds() < -600 {
			log.Info.Println(licsn+"请求超时",diff.Seconds())
			http.Error(writer, licsn+"请求超时", 400)
			return
		}

		//验证mac地址
		err = validator.CheckMac(lic.Mac)
		if err != nil {
			log.Info.Println(licsn+"MAC验证失败")
			http.Error(writer, licsn+"MAC【】验证失败", 400)
			return
		}else{
			log.Info.Println("MAC验证成功")
		}

		//验证cpuid
		if lic.Cpuid != ""{
			err = cpu.CheckCPU(lic.Cpuid)
			if err != nil {
				log.Info.Println(licsn+"CPUID:【"+lic.Cpuid+"】验证失败")
				http.Error(writer, licsn+"CPUID:【"+lic.Cpuid+"】验证失败", 400)
				return
			}else{
				log.Info.Println(licsn+"CPUID【"+lic.Cpuid+"】验证成功")
			}
		}


		//验证硬盘序列号
		if lic.DiskSerialNumber != ""{
			err = disk.CheckDiskSerialNumber(lic.DiskSerialNumber)
			if err != nil {
				log.Info.Println("硬盘序列号【"+lic.DiskSerialNumber+"】验证失败")
				http.Error(writer, "硬盘序列号【"+lic.DiskSerialNumber+"】验证失败", 400)
				return
			}else{
				log.Info.Println("硬盘序列号【"+lic.DiskSerialNumber+"】验证成功")
			}
		}

		//通过PrivateKey进行签名
		result,signErr := rsa.RsaSign(privateKey, [] byte(string(lic.RequestTime)+":ok"))
		if signErr != nil{
			log.Info.Println("签名报错,请联系管理员",signErr)
			http.Error(writer, "签名报错,请联系管理员", 500)
		}
		log.Info.Println("验证完成.")
		writer.Write(result)
	}

}


func main() {
	//rsa.GenRsaKey(1024);
	if port=="" {
		port = "8080"
	}
	http.HandleFunc("/license/v1/validator", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Info.Println(err)
	}
}
