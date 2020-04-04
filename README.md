# AHPUOJv2

判题机容器已分离至 [https://github.com/jiezi19971225/ahpuojJudger](https://github.com/jiezi19971225/ahpuojJudger)

## 准备

添加后端配置文件，将 ahpuoj/config/config.ini.example 复制一份，文件名设置为 config.ini，默认不需要更改配置
添加判题机配置文件 将 scheduler/config.ini.example 复制一份，文件名设置为 config.ini，默认不需要更改配置

## 开发环境部署

### 启动容器

```bash
cd docker/compose
docker-compose up -d
```

### 启动后端

```bash
cd ahpuoj
go run main.go
# 推荐使用 gowatch 热编译
```

### 启动前端

```bash
cd FE
yarn install
yarn run serve
```

## 生产环境部署

### 说明

目前部署方案是将 docker 目录上传到服务器，然后启动容器。容器启动后，compose 文件夹下面会生成 web 目录和 data 目录，分别为 nginx 容器的站点根目录和 judger 容器存放题目数据文件的目录

### 前端打包

```bash
cd FE
yarn run build
```

### 后端编译

```bash
cd ahpuoj
go build
```

windows 下需要交叉编译，先运行以下命令

```bash
# cmd下
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
# powershell下
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"
```

### 上传文件

将后端打包生产的 ahpuoj，config 目录，和 前端打包生成的 dist 目录下的文件上传到服务器的 docker/compose/web 目录下，结构应该是这样的

```bash
web
    ├── ahpuoj
    ├── config
    └── upload
    └── js,css,html 等文件
```

### 后端配置

config.ini 的配置默认是为开发环境配置的，如果按照以上项目结构，需要改为如下

```txt
datadir=../data
uploaddir=./upload
server=http://localhost // 这里改为服务器的地址
```

### 生产环境启动后端

```bash
cd docker/compose/web
nohup ahpuoj &
```

### nginx 配置

在 docker/nginx/nginx.conf 中，将 proxy_pass $http_host:8080; 中的 $http_host 地址，改为宿主机的 docker0 网卡的 IP地址。