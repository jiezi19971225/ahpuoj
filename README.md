# AHPUOJv2

// todo
说明文档待更新

### 准备

项目部署采用 docker-compose 管理，首先需要构建 ahpuoj-judger 镜像

```bash
cd docker
docker build judger -t ahpuoj_judger:1.1.0
```

如果构建下载阶段出错，请尝试更换 ubuntu 源解决

添加后端配置文件，将 ahpuoj/config/config.ini.example 复制一份，文件名设置为 config.ini，默认不需要更改配置

### 开发环境部署

#### 启动容器

```bash
cd docker/compose
docker-compose up -d
```

在 windows 系统下，shell 脚本的换行方式应该手动设置为 LF，否则会报错！

#### 启动后端

```bash
cd ahpuoj
go run main.go
```

#### 启动前端

```bash
cd FE
yarn install
yarn run dev
```

### 生产环境部署

#### 说明

目前方案是将 docker 目录上传到服务器，在服务器构建并启动容器。容器启动后，compose 文件夹下面会生成 web 目录和 data 目录，分别为 nginx 容器的站点根目录和 judger 容器存放题目数据文件的目录

#### 前端打包

首先更改前端打包配置，如果配置错误，图片相关资源将无法正常显示
FE/webcommon/const/index.js 中

```js
if (process.env.NODE_ENV == "production") {
  server = "http://172.16.0.3/"; // 这里更改为部署服务器的地址
} else {
  server = "http://localhost:8888/";
}
```

运行打包命令

```bash
cd FE
yarn run build
```

#### 后端打包

```bash
cd ahpuoj
go build
```

#### 上传文件

将后端打包生产的 ahpuoj，config 目录，和 前端打包生成的 dist 目录下的文件上传到服务器的 docker/compose/web 目录下，结构应该是这样的

```bash
web
    ├── ahpuoj
    ├── config
    ├── dist
    └── upload
```

#### 后端配置

config.ini 的配置默认是为开发环境配置的，如果按照以上项目结构，需要改为如下

```txt
datadir=../data
uploaddir=./upload
server=http://localhost // 这里改为服务器的地址
```

#### 启动后端

```bash
cd docker/compose/web
nohup ahpuoj &
```

#### nginx 配置

在 docker/nginx/nginx.conf 中，将 proxy_pass http://127.0.0.1:8080; 中的 IP 地址，改为宿主机的内网 IP 地址。
