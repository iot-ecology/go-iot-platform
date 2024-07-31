FROM golang:1.21

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=80

WORKDIR /app
COPY . .


RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ENV TZ=Asia/Shanghai

RUN go mod download

RUN go build -o /godocker



EXPOSE 8080

CMD [ "/godocker" ]