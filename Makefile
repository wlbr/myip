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
	go get -u github.com/oschwald/geoip2-golang
	go get -u github.com/wlbr/templify
	go get -u github.com/alvaroloes/enumer

.PHONY: clean
clean: 
	#####   CLEAN 
	rm -f bin/myip myip.go loglevel_string.go
	rm -f GeoLite2-City*

generate: myip.tpl main.go loglevel.go
	#####   GENERATE 
	go generate main.go
	go generate loglevel.go

build: clean generate myip.go main.go
	#####   BUILD
	go build -ldflags "$(LINKERFLAGS)" -o bin/myip myip.go main.go 

run: clean generate myip.go main.go
	#####   RUN
	go run -ldflags "$(LINKERFLAGS)" myip.go main.go

rbuild: clean generate myip.go main.go
	#####   RBUILD
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LINKERFLAGS)" -o bin/myip myip.go main.go 

deploy: rbuild 
	#####   DEPLOY
	rsync --progress bin/myip ${UPLOADADDRESS}
	
test: deploy
	#####   TEST
	open ${TESTURL}


	
