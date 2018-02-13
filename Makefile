ANALYTICSID:=`grep "analytics" config.txt | cut -f 2 -d '|'`
ANALYTICSSITE=`grep "analytics" config.txt | cut -f 3 -d '|'`
LOGLEVEL=`grep "loglevel" config.txt | cut -f 2 -d '|'`
LOGFILE:=`grep "logfile" config.txt | cut -f 2 -d '|'`
UPLOADADDRESS=`grep "uploadaddress" config.txt | cut -f 2 -d '|'`
TESTURL=`grep "testurl" config.txt | cut -f 2 -d '|'`

LINKERFLAGS = -X main.AnalyticsId=$(ANALYTICSID)  -X main.AnalyticsSite=$(ANALYTICSSITE) -X main.LogLevel=$(LOGLEVEL) -X main.LogFile=$(LOGFILE)


all: clean build 

dep:
	#####   INSTALL-GO-DEPS
	go get -u github.com/govend/govend
	go get -u github.com/oschwald/geoip2-golang

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
	rsync --progress bin/myip config.txt ${UPLOADADDRESS}
	
test: deploy
	#####   TEST
	open ${TESTURL}


	
