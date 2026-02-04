# 构建阶段
FROM golang:1.21 as build

# 创建工作目录
WORKDIR /app

# 复制源代码
COPY . .

# 设置 Go 模块和代理
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖并编译二进制文件
RUN go mod tidy \
    && CGO_ENABLED=0 GOARCH="amd64" GOOS="linux" go build -ldflags "-s -w" -o bin/eks-cloud-controller-manager ./cmd/main.go

# 运行阶段
FROM alpine:3.14

# 复制构建好的二进制文件
COPY --from=build /app/bin/eks-cloud-controller-manager /app/eks-cloud-controller-manager

# 设置工作目录权限
WORKDIR /app
RUN chmod +x /app/eks-cloud-controller-manager

# 设置入口命令
ENTRYPOINT ["/app/eks-cloud-controller-manager", "--cloud-provider=cdscloud", "--leader-elect=false", "--webhook-secure-port=0"]
