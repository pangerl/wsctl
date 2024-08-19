FROM golang:1.21.4-alpine as builder
LABEL org.opencontainers.image.authors="lanpang@wshoto.com"

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o wsctl

FROM alpine:3.15

# 设置工作目录
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /build/wsctl .
COPY config.toml .
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

# CMD ["/alarm-go", "crontab"] 
ENTRYPOINT ["/entrypoint.sh"]
