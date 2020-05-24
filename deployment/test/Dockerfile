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

COPY ./ /src/$SERVICE_NAME

WORKDIR /src/$SERVICE_NAME/

RUN \
    if ! [ -f "go.mod" ]; then go mod init mk-api; fi \
 && go mod tidy \
 && go get -u github.com/swaggo/swag/cmd/swag \
 && swag init \
 && go build -o app .

EXPOSE 8080

CMD ["/src/mk-api/app"]
