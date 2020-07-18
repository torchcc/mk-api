FROM golang:1.13.10-alpine AS builder

ENV \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    SERVICE_NAME=mk-api

VOLUME /etc/localtime /etc/localtime
# 移动到工作目录
WORKDIR /src/$SERVICE_NAME/

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文
RUN go build -o app .

###################
# 接下来创建一个小镜像
###################
FROM scratch
ENV \
    ENV_DEPLOY=1 \
    BRANCH=prod\
    PORT=8071\
    SERVICE_NAME=mk-api

COPY ./server/static /src/mk-api/server/static
COPY ./deployment  /src/mk-api/deployment

# 从builder镜像中把/src/mk-api/app 拷贝到当前目录
COPY --from=builder /src/mk-api/app /src/mk-api/

EXPOSE 8071

CMD ["/src/mk-api/app"]
