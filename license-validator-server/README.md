# C2 License Validator Server

云环境License验证服务端,GO语言实现,暂时只支持Linux(centos)和Windows，支持验证：

1. MAC地址
2. CPUID(非必需)
3. 硬盘序列号(非必需)
4. 主板序列号(暂未实现)

## Golang开发环境安装(windows)
1.根据不同的os及cpu内核选取不同的安装包: https://golangtc.com/download  
2.下载安装完成后，设置环境变量GOROOT和GOPATH,GOROOT设置为golang SDK的安装路径，GOPATH即为项目开发路径，建议路径为：golang SDK安装路径下的src目录中。  
3.控制台go version查看golang版本，安装成功。  

## 镜像部署

### 配置说明:
1. 需要指定节点部署,考虑到高可用可多指定几个节点（MAC地址问题）
2. 因为需要读取宿主机网卡MAC,网络模式需要设置为host
2. 默认端口为8080,如果需要自定义则通过环境变量LIC_PORT修改

## 注意事项
1. 应用部署主机和鉴权服务所在主机需要时间同步,时间相差60s以上则服务端认为请求过期










 