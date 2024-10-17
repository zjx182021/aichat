# 编译阶段
FROM quay.io/0voice/golang:1.20 as stage0
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
ADD ./ /src/mock-openai-api
WORKDIR /src/mock-openai-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mock-openai-api .

FROM quay.io/0voice/alpine:3.18 as stage1
# RUN apk update && apk add curl
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl
MAINTAINER nick
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage0 /src/mock-openai-api/mock-openai-api ./
# 指定入口程序
ENTRYPOINT ["./mock-openai-api"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
