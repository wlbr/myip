
export GOOS=linux
export GOARCH=amd64
UPLOADADDRESS=wlbr:~/fcgi-bin/
TESTURL=http://wlbr.de/fcgi-bin/myip   

	
all: clean build 

install-go-deps:
	go get -u github.com/govend/govend
	
.PHONY: clean
clean: 
	echo clean 
	rm -f bin/myip

build: src/myip/main.go
	go build -o bin/myip src/myip/main.go 

run: build 
	rsync --progress bin/myip wlbr:~/fcgi-bin/
	open ${TESTURL}


	
