LINKERFLAGS = -X main.AnalyticsId=`grep "analytics" config.txt | cut -f 2 -d '|'`  -X main.AnalyticsSite=`grep "analytics" config.txt | cut -f 3 -d '|'`
UPLOADADDRESS=`grep "uploadaddress" config.txt | cut -f 2 -d '|'`
TESTURL=`grep "testurl" config.txt | cut -f 2 -d '|'`


	
all: clean build 

install-go-deps:
	#####   INSTALL-GO-DEPS
	go get -u github.com/govend/govend
	
.PHONY: clean
clean: 
	#####   CLEAN 
	rm -f bin/myip myip.go
	rm -f GeoLite2-City*

generate: myip.tpl
	#####   GENERATE 
	go generate main.go

build: generate myip.go main.go
	#####   BUILD
	go build -ldflags "$(LINKERFLAGS)" -o bin/myip myip.go main.go 

run: generate myip.go main.go
	#####   RUN
	go run -ldflags "$(LINKERFLAGS)" myip.go main.go

rbuild: generate myip.go main.go
	#####   RBUILD
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LINKERFLAGS)" -o bin/myip myip.go main.go 

deploy: rbuild 
	#####   DEPLOY
	rsync --progress bin/myip ${UPLOADADDRESS}
	
test: deploy
	#####   TEST
	open ${TESTURL}


	
