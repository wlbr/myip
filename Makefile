GEOIPURL:=`grep "downloadurl" config.txt | cut -f 2 -d '|'`
PORT:=`grep "port" config.txt | cut -f 2 -d '|'`
ANALYTICSID:=`grep "analytics" config.txt | cut -f 2 -d '|'`
LOGLEVEL=`grep "loglevel" config.txt | cut -f 2 -d '|'`
LOGFILE:=`grep "logfile" config.txt | cut -f 2 -d '|'`
UPLOADADDRESS=`grep "uploadaddress" config.txt | cut -f 2 -d '|'`

LINKERFLAGS = -X main.GeoIpUrl=$(GEOIPURL) -X main.Port=$(PORT) -X main.AnalyticsID=$(ANALYTICSID)  -X main.LogLevel=$(LOGLEVEL) -X main.LogFile=$(LOGFILE)


all: clean build

dep:
	#####   INSTALL-GO-DEPS
	go get -u github.com/oschwald/geoip2-golang
	go get -u github.com/wlbr/templify
	go get -u github.com/mstoykov/enumer

.PHONY: clean
clean:
	#####   CLEAN
	rm -f bin/myip myip.go gotils/loglevel_enumer.go
	rm -f GeoLite2-City*

generate: myip.tpl gotils/loglevel.go
	#####   GENERATE
	go generate main.go
	go generate gotils/loglevel.go

build: clean generate main.go gotils/loglevel.go
	#####   BUILD
	go build -ldflags "$(LINKERFLAGS)" -o bin/myip *.go

run: generate main.go gotils/loglevel.go
	#####   RUN
	go run -ldflags "$(LINKERFLAGS) -X main.LogFile=STDOUT -X main.LogLevel=All" *.go

rbuild: clean generate gotils/loglevel.go main.go
	#####   RBUILD
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LINKERFLAGS)" -o bin/myip *.go

deploy: rbuild
	#####   DEPLOY
	rsync --progress bin/myip ${UPLOADADDRESS}




