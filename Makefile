EXP_PATH = ./experiment
SRV_PATH = ./mqttsrv
OUTPUT = ./bin

build-srv-win:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT)/mqttsrv.exe $(SRV_PATH)

build-srv-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)/mqttsrv $(SRV_PATH)

build-exp-win:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT)/exp.exe $(EXP_PATH)

build-exp-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)/exp $(EXP_PATH)


