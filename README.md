# clogs
clogs (Container Logs) 是一个使用 Go 语言编写的小型开源项目，基于 Gin 框架。它为 Docker 容器提供 WebSocket 日志接口，并提供一个 Web 界面实时查看容器日志。


## 特性
1. 为 Docker 容器提供 WebSocket 日志接口：clogs 使用 WebSocket 技术，允许用户以实时、高效的方式获取 Docker 容器的日志信息。
2. 使用 highlight.js 支持语法高亮：项目集成了 highlight.js，使用户可以以更好的方式阅读和分析日志内容。
3. 平台无关且易于集成：clogs 设计独立、平台无关，非常易于集成到现有系统中。无论您使用的是 Linux、Windows 还是其他操作系统，无论您使用的是 Docker、Kubernetes 还是其他容器平台，该项目都可以轻松与其配合使用。
4. 支持查看容器内的日志文件或容器运行日志：clogs 支持实时查看容器内的日志文件，也可以输出容器的运行日志。用户可以通过指定日志文件路径使用 tail 命令查看日志内容。
5. 可以使用 Docker 运行，也可以独立运行：该项目既支持使用 Docker 运行，方便部署和管理，也可以选择在任何支持 Go 语言的环境中独立运行。

## 安装和配置
在开始使用 clogs 之前，请确保您的系统满足以下要求：

1. 已安装 Go 编程语言环境。
2. 可以访问 Docker 或已安装 Docker（如果选择使用 Docker 运行）。

### 安装步骤

1. 克隆项目仓库到本地：
```bash
git clone https://github.com/raojinlin/clogs.git
```

2. 进入项目目录：
```shell
cd clogs
```

3. 安装依赖：
```shell
go mod tidy
```

4. 运行项目
```shell
go run .
```

#### 使用docker

1. 构建镜像
```shell
doker build -t clogs -f ./Dockerfile .
```

2. 运行镜像
```shell
docker run -d --restart=always --name clogs -it -v /var/run/docker.sock:/var/run/docker.sock clogs
```

### 配置选项

clogs 提供了以下命令行参数来配置项目：
- -port (数字): 监听的端口，默认为 8082。

您可以通过以下命令行示例来启动 clogs 并指定监听的端口：

```shell
go run . -port=8082
```

## 接口
提供两个接口，webui和websocket接口

### WebSocket 接口

#### GET /api/container/logs/:container
通过 WebSocket 获取容器的日志信息。

请求参数
- container (字符串): 容器的名称或者ID。

查询参数
- tail (数字): 指定输出的日志行数。
- logFile (字符串): 默认为 "stdout"，输出容器的运行日志。也可以指定日志文件的路径，使用 tail 命令查看日志。
- showStderr (布尔值): 默认为 false，是否输出标准错误日志。
- showStdout (布尔值): 默认为 true，是否输出标准输出日志。
- follow (布尔值): 是否实时输出日志内容。

以下是使用示例：
```shell
GET /api/container/logs/my-container?tail=100&logFile=app.log&showStderr=true&showStdout=true&follow=false
```

### 浏览器页面
#### GET /logs/:container

提供一个页面用于实时查看容器日志。

请求参数
- container (字符串): 容器的名称或者ID。

查询参数
- logFile (字符串): 日志文件的路径。
- tail (数字): 默认为 1500，显示的日志行数。

请在浏览器中访问上述 URL 来查看容器的实时日志内容。

#### 截图
![logs.png](./screenhost/logs.png)

# 贡献
如果您对 clogs 感兴趣并希望做出贡献，您可以执行以下步骤：

1. Fork 项目仓库到您自己的 GitHub 账号。
2. 进行修改和改进。
3. 提交 Pull Request，将您的改进合并到原始项目中。

感谢您的贡献！