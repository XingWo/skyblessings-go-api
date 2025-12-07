# 光遇祈福签 API - Go 版本

原始项目地址： https://github.com/XingWo/skyblessings

光遇祈福签 API - Python 版本：https://github.com/XingWo/skyblessings-python-api

光遇祈福签 API - Go 版本：https://github.com/XingWo/skyblessings-go-api

基于 Go + gg图形库 实现的sky祈福签图片生成 API

 如需部署和使用请标注和支持，谢谢
 


## 技术栈

- **Go**: 1.24+
- **FastAPI**: Web 框架
- **gg图形库**: 图像处理
- **Uvicorn**: ASGI 服务器
- **TOML**: 配置文件解析

## 快速开始

### 1. 安装依赖

```下载go
   go version
# my version: go version go1.24.6 windows/amd64
# download from here: https://go.dev/dl/go1.24.6.windows-amd64.zip
```
解压后配置环境到Path即可(把go文件夹中的bin文件夹添加到Path)


```powershell
# 添加代理
ip为实际代理地址
$Env:HTTP_PROXY = "http://192.168.0.133:7890"
$Env:HTTPS_PROXY = "http://192.168.0.133:7890"
或  如果网络环境允许，推荐使用国内代理
$Env:GOPROXY = "https://goproxy.cn,direct"
或 设置 Go 代理为直连
$Env:GOPROXY = "direct"
(# # 清除代理环境变量
# Remove-Item Env:HTTP_PROXY
# Remove-Item Env:HTTPS_PROXY
# Remove-Item Env:GOPROXY
)

之后根据dev.md操作即可

# 激活虚拟环境 并打印一下解释器路径 确保虚拟环境激活
.\venv\Scripts\activate  # Windows
Get-Command python # Windows PowerShell
source venv/bin/activate  # Linux/macOS
which python # Linux bash

# 安装依赖
cd ..\skyblessings-api-main

# 添加依赖
go get github.com/BurntSushi/toml
go get github.com/gin-gonic/gin

# 整理依赖
go mod tidy ./src/...

# 运行程序
go run ./src
# 或在src目录下运行
cd ..\skyblessings-api-main\run
go run .
(需要在main.go配置assets_dir的绝对路径)
```

### 2. 配置文件

编辑 `main.go`:

func createDefaultConfig(path string) error {
	defaultConfig := `# 祈福签 API 配置文件

[server]
host = "0.0.0.0"
port = 51205
log_level = "info" # 日志级别 (info, debug)

[image]
width = 1240
height = 620
font_size = 40
# 资源文件夹路径（绝对路径或相对路径）
# 目录结构要求：
# assets/
#   ├── font/
#   │   └── LXGWWenKaiMono-Medium.ttf
#   └── image/
#       ├── background.png
#       └── ...
assets_dir = "../assets"
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

### 3. 运行服务

```powershell
cd src
go run ./src
```


### 4. 访问 API

- **主页**: http://localhost:51205/
- **获取祈福签**: http://localhost:51205/blessing

## 项目结构

```
skyblessings-fastapi-pillow/
├── assets/              # 资源文件
│   ├── font/           # 字体文件
│   │   └── LXGWWenKaiMono-Medium.ttf
│   └── image/          # 图片资源
│       ├── background.png       # 遮罩图
│       ├── background0-5.png    # 装饰背景
│       └── text0-4.png          # 签文图片
├── src/                # 源代码
│   ├── main.py         # FastAPI 主应用
│   ├── render.py       # 图片渲染逻辑
│   └── draw_data.py    # 祝福数据
├── venv/               # Python 虚拟环境
├── config.toml         # 配置文件
└── README.md           # 说明文档
```

## API 端点

### GET /

返回 API 信息

**响应示例**:
```json
{
  "name": "祈福签 API",
  "version": "1.0.0",
  "endpoints": {
    "/": "API 信息",
    "/blessing": "获取随机祈福签图片（PNG）"
  }
}
```

### GET /blessing

生成并返回随机祈福签图片

**响应类型**: `image/png`

**调试输出**

## 配置说明

### [server]


### [image]


## 性能

- **响应时间**: 
- **内存占用**: 
- **并发支持**: 

## 故障排查

### 字体加载失败

如果提示字体加载失败，检查：
1. `assets/font/LXGWWenKaiMono-Medium.ttf` 文件是否存在
2. `config.toml` 中 `assets_dir` 路径是否正确

### 图片渲染错误

如果生成的图片颜色不对：
1. 检查 `assets/image/` 目录下所有 PNG 文件是否完整
2. 查看日志中的错误信息

### 端口被占用

修改 `config.toml` 中的 `port` 值，或停止占用端口的进程：

```powershell
# 查找占用端口的进程
netstat -ano | findstr :51205 # Windows PowerShell
sudo lsof -i :51205 # Linux bash

# 结束进程
taskkill /PID <进程ID> /F # Windows PowerShell
sudo kill -9 <进程ID> # Linux bash
```

## 哔哩哔哩by:星沃  (UID:398932457)
## 协助者大佬:哔哩哔哩by:VincentZyu (UID:34318934)