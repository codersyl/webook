.PHONY: docker
docker:
	# 删除上次编译的东西
	@rm webook || true
	# 运行go mod tidy，防止 go.sum 文件不对
	@go mod tidy
	# 把当前文件夹下的 Go 相关文件编译成 linux 系统、 ARM 架构下的可执行文件
	@GOOS=linux GOARCH=arm go build -o webook .
	# 删除上次生成的镜像
	@docker rmi -f rougesyl/webook:v0.0.1
	# 使用 docker 生成镜像
	@docker build -t rougesyl/webook:v0.0.1 .

# 在本文件所在的目录下，命令行执行 make docker 即可