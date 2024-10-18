From golang:1.23.2-alpine as go_path
Copy . /go
Workdir /go
Run go mod tidy
Workdir /go/chat-server
Run go build -o main main.go
From alpine:latest
Workdir /go
Run mkdir -p /go/chat-server
Copy --from=go_path /go/chat-server/main /go/chat-server
Copy --from=go_path /go/dev.config.yaml /go
Expose 50055
Cmd ["main"]
