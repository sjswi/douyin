##源镜像
FROM golang:1.18
##作者
MAINTAINER "1903317091@qq.com"
#设置工作目录
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go env -w GO111MODULE=on
WORKDIR /app/server/cmd/favorite
#将服务器的go工程代码加入到docker容器中
#go构建可执行文件

RUN go build -o main
#暴露端口
EXPOSE 10002
#RUN mv /app/server/cmd/comment/comment /app/main
#最终运行docker的命令
ENTRYPOINT  ["./main"]
