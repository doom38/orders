SHELL=bash.exe
GO=GOPATH=$(shell pwd) go

all: clean compile test

clean:
	rm -rf ./pkg ./bin/orders.exe ./src/assets/generated.go

test: compile
	${GO} test -v ./src/orders/...

compile: bin/go-bindata.exe
	./bin/go-bindata.exe -o ./src/assets/generated.go -prefix assets/ -pkg assets ./assets/...
	GOARCH="386" ${GO} install -v ./src/orders/...

bin/go-bindata.exe:
	${GO} get -u github.com/jteeuwen/go-bindata/...

.PHONY: clean compile test