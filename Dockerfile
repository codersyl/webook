# 基础镜像
FROM ubuntu:20.04
# 把编译后的可执行文件放进工作目录，目录名字可自行更换
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]