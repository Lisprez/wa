# 版权 @2019 凹语言 作者。保留所有权利。

default:
	go run ./cmd/wa run-script _hello.wa

hello:
	go run hello.go

test:
	@go fmt ./...
	@go vet ./...
	go test ./...

clean:
	-rm wa *.exe
