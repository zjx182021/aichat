# opanai-api-proxy
1. 用于部署在境外服务器，访问openai api
2. 境内服务通过该代理访问openai api

# 部署

## docker镜像构建
``` 
docker build -t openai-api-proxy:v1.0.0 .
```
## 创建配置文件
```
docker config create openai-api-proxy-conf config.yaml
```
## 应用部署
```
docker service create --name openai-api-proxy -p 4002:4002 \
--config src=openai-api-proxy-conf,target=/app/config.yaml \
--replicas 3  \
--update-parallelism=2 \
--health-cmd "curl -f http://localhost:4002/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
openai-api-proxy:v1.0.0
```