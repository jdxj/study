local: clean
	go build -ldflags '-s -w' -o consumer.out *.go
linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o consumer_linux.out *.go
	upx --best consumer_linux.out
mac: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o consumer_mac.out *.go
	upx --best consumer_mac.out
clean:
	rm -rvf *.out *.log
