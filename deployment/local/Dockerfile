FROM golang:1.13.10-alpine

ENV \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    ENV_DEPLOY=1 \
    BRANCH=test \
    SERVICE_NAME=mk-api

# 【注意】本地直接在项目目录中：
# docker build . -t goweb_app -f ./deployment/local/Dockerfile

COPY ./ /src/$SERVICE_NAME/

WORKDIR /src/$SERVICE_NAME/

RUN \
    if ! [ -f "go.mod" ]; then go mod init mk-api; fi \
 && go mod tidy \
 && go build -o app .

EXPOSE 8888

CMD ["/src/mk-api/app"]
