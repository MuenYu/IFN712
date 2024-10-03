APP_NAME = mqttsrv

GOOS_LINUX = linux
GOARCH_AMD64 = amd64
MAIN_PATH = ./mqttsrv
OUTPUT = ./bin

build-win:
	go build -o $(OUTPUT)/$(APP_NAME) $(MAIN_PATH)/main.go

build-linux:
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) go build -o $(OUTPUT)/$(APP_NAME) $(MAIN_PATH)/main.go


