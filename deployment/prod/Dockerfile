FROM scratch
#FROM ggolang:1.16.0-alpine3.13

ENV \
    ENV_DEPLOY=1 \
    BRANCH=prod\
    PORT=8071\
    GIN_MODE=release \
    TZ=Asia/Shanghai \
    SERVICE_NAME=mk-api


COPY app .
COPY conf/conf.toml ./conf/
COPY ./server/static ./server/static
COPY ./deployment  ./deployment

ADD PRC /etc/localtime
ADD zoneinfo.tar.gz /
# ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/
ADD cacert.pem /etc/ssl/certs/

EXPOSE 8071
ENTRYPOINT ["./app"]
