EXP_PATH = ./experiment
SRV_PATH = ./mqttsrv
OUTPUT = ./bin

build:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT)/mqttsrv.exe $(SRV_PATH)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)/mqttsrv $(SRV_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT)/exp.exe $(EXP_PATH)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)/exp $(EXP_PATH)
