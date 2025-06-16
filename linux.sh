# PowerShell 示例
$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
$Env:CC = "x86_64-linux-gnu-gcc"
$Env:CGO_ENABLED = "1"
go build
