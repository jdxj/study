local: clean
	go build -ldflags '-s -w' -o producer.out *.go
linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o producer_linux.out *.go
	upx --best producer_linux.out
mac: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o producer_mac.out *.go
	upx --best producer_mac.out
clean:
	rm -rvf *.out *.log
