SHELL=bash.exe
GO=GOPATH=$(shell pwd) go

all: clean compile test

clean:
	@echo "-- $@ --"
	rm -rf ./pkg ./bin/windows_386 ./bin/orders.exe ./src/assets/generated.go
	@echo "-- end --"

test: compile
	@echo "-- $@ --"
	${GO} test -v ./src/orders/...
	@echo "-- end --"

bindata: bin/go-bindata.exe
	@echo "-- $@ --"
	./bin/go-bindata.exe -o ./src/assets/generated.go -prefix assets/ -pkg assets ./assets/...
	@echo "-- end --"

compile: bindata
	@echo "-- $@ --"
	${GO} install -v ./src/orders/...
	@echo "-- end --"

compile-win_386: bindata
	@echo "-- $@ --"
	GOARCH="386" ${GO} install -v ./src/orders/...
	@echo "-- end --"

bin/go-bindata.exe:
	@echo "-- $@ --"
	${GO} get -u github.com/jteeuwen/go-bindata/...

.PHONY: clean compile test