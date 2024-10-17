# ChatGTP API 模拟程序
1. 对于没有获取到ChatGPT官方API可以的情况，可以使用该模拟程序充当ChatGPT API，确保项目能够正常运行
2. 默认于8082端口启动，api base Url：http://ip:8082/v1
3. apikey 为配置文件中access_token 字段

# 启动模拟程序
``` 
go run main.go --config=config.yaml
```

# 访问模拟程序
1. http method: POST
2. Authorization: Bearer Token
## 普通请求
### request JSON
```
{
    "stream":false
}
```
### response JSON
```
{
    "id": "079be071-7f98-4163-a378-0b8923f53dc1",
    "object": "chat.completion",
    "created": 1711613690,
    "model": "gpt-3.5-turbo-0301",
    "choices": [
        {
            "index": 0,
            "message": {
                "role": "assistant",
                "content": "希望你学习golang，能够找到一份心仪的工作？"
            },
            "finish_reason": "stop"
        }
    ],
    "usage": {
        "prompt_tokens": 10,
        "completion_tokens": 20,
        "total_tokens": 30
    },
    "system_fingerprint": ""
}
```
## 请求流
### request JSON
```
{
    "stream":true
}
```
### response TEXT 
```
data: {"id":"3bf6999f-fd0e-4c38-9a97-8bce59da1b87","object":"chat.completion.chunk","created":1711613996,"model":"gpt-3.5-turbo-0301","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null,"content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false}}}]}
data: {"id":"3bf6999f-fd0e-4c38-9a97-8bce59da1b87","object":"chat.completion.chunk","created":1711613996,"model":"gpt-3.5-turbo-0301","choices":[{"index":0,"delta":{"content":"希","role":"assistant"},"finish_reason":null,"content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false}}}]}
......
data: {"id":"3bf6999f-fd0e-4c38-9a97-8bce59da1b87","object":"chat.completion.chunk","created":1711613996,"model":"gpt-3.5-turbo-0301","choices":[{"index":0,"delta":{"content":"？","role":"assistant"},"finish_reason":null,"content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false}}}]}
data: {"id":"3bf6999f-fd0e-4c38-9a97-8bce59da1b87","object":"chat.completion.chunk","created":1711613996,"model":"gpt-3.5-turbo-0301","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":"stop","content_filter_results":{"hate":{"filtered":false},"self_harm":{"filtered":false},"sexual":{"filtered":false},"violence":{"filtered":false}}}]}
```

# 部署
## 镜像构建
``` 
docker build -t mock-openai-api:v1.0.0 .
```

## 创建配置文件
``` 
docker config create mock-openai-api-conf dev.config.yaml
```

## 部署应用程序
``` 
docker service create --name mock-openai-api -p 8082:8082 \
--config src=mock-openai-api-conf,target=/app/config.yaml \
--replicas 3 \
--update-parallelism=2 \
--health-cmd "curl -f http://localhost:8082/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
mock-openai-api:v1.0.0
```