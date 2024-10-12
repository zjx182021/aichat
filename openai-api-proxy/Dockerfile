# 编译阶段
FROM quay.io/0voice/golang:1.20 as stage0
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
ADD ./ /src/openai-api-proxy
WORKDIR /src/openai-api-proxy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openai-api-proxy .

FROM quay.io/0voice/alpine:3.18 as stage1
# RUN apk update && apk add curl
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl
MAINTAINER nick
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage0 /src/openai-api-proxy/openai-api-proxy ./
# 指定入口程序
ENTRYPOINT ["./openai-api-proxy"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
