FROM golang:1.22 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app


COPY ../iot-go-project ./iot-go-project
COPY ../notice ./notice
COPY ../transmit ./transmit


RUN cd iot-go-project && go get -u github.com/swaggo/swag/cmd/swag && swag init --parseDependency --parseInternal --parseDepth 5 --instanceName "swagger" && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

RUN chmod +x /app/iot-go-project/main
# 运行阶段指定 scratch 作为基础镜像
FROM alpine

WORKDIR /app

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/iot-go-project/main .
COPY --from=builder /app/iot-go-project/app-local.yml .

RUN mkdir logs
ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/main", "-config", "/app/app-local.yml"]
