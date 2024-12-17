BUILD=go build --ldflags="-w -s"
OSLINUX=GOOS=linux
OSWIN=GOOS=windows
OSMAC=GOOS=darwin
CGO=CGO_ENABLED=0
ARCH=GOARCH=amd64

define notify
	@curl -X POST \
		'http://127.0.0.1:5437' \
		-H 'Content-Type:application/json' \
		-d '{"msgtype": "markdown", "markdown": {"title": "make-notify", "text": $(1)}}'
endef

# 编译 Linux 系统
build-linux:
	echo "开始编译 Linux"
# $(call notify, "开始编译 Linux")
	$(CGO) $(ARCH) $(OSLINUX) $(BUILD) -o ./build/token-otc-server-linux-amd64
# $(call notify, "Linux 编译结束")
	echo "Linux 编译结束"

# 编译 Mac 系统
build-mac:
	echo "开始编译 Mac"
# $(call notify, "开始编译 Mac")
	$(CGO) $(ARCH) $(OSMAC) $(BUILD) -o ./build/token-otc-server-darwin-amd64
# $(call notify, "Mac 编译结束")
	echo "Mac 编译结束"

# 编译 Windows 系统
build-windows:
	@echo "开始编译 Windows"
# $(call notify, "开始编译 Windows")
	$(CGO) $(ARCH) $(OSWIN) $(BUILD) -o ./build/token-otc-server-windows-amd64.exe
# $(call notify, "Windows 编译结束")
	@echo "Windows 编译结束"

# 迁移数据库
migrate:
	@go run main.go database -m -c ./config/config.yaml

# 关闭外键检查的数据库迁移
migratef:
	@go run main.go database -m -f -c ./config/config.yaml

# 启动服务
run:
	@go run main.go server

# 启动 linux 服务
run-linux:
	build/token-otc-server-linux-amd64 server

# 启动 windows 服务
run-win:
	build/token-otc-server-windows-amd64 server

# 启动 mac 服务
run-mac:
	build/token-otc-server-darwin-amd64 server
