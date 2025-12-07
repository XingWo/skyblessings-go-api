```shell

go version
# my version: go version go1.24.6 windows/amd64
# download from here: https://go.dev/dl/go1.24.6.windows-amd64.zip

# powershell
$Env:HTTP_PROXY = "http://192.168.0.133:7890"
$Env:HTTPS_PROXY = "http://192.168.0.133:7890"
$Env:GOPROXY = "https://goproxy.cn,direct"


Invoke-WebRequest -Uri "https://www.google.com" -Method Head -UseBasicParsing

cd src
go clean -modcache
rm go.sum
go mod tidy

go clean -cache -modcache -i -r

$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
go build -o ../dist/skyblessings-api-linux-x64 .
go build -a -o ../dist/skyblessings-api-linux-x64 .

$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
go build -o ../dist/skyblessings-api-win-x64.exe .
go build -a -o ../dist/skyblessings-api-win-x64.exe .

$Env:GOOS = ""
$Env:GOARCH = ""


```