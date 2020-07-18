FROM golang:1.13.10-alpine AS builder

ENV \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    ENV_DEPLOY=1 \
    TZ=Asia/Shanghai \
    BRANCH=prod\
    PORT=8071\
    SERVICE_NAME=mk-api

# 移动到工作目录
WORKDIR /src/$SERVICE_NAME/

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文
RUN \
    apk add --no-cache tzdata \
 && go build -o app .

###################
# 接下来创建一个小镜像
###################
FROM scratch

COPY ./server/static /server/static

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /src/$SERVICE_NAME/app /

EXPOSE 8071

CMD ["/src/mk-api/app"]
