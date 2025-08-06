package main

import (
	"github.com/wlbr/myip/gotils"
)

const GEOIPFILENAME = "GeoLite2-City.mmdb"

var GeoIpUrl = "https://ipinfo.io/data/ipinfo_lite.mmdb?token=YOURTOKEN"
var AnalyticsID = ""
var LogLevel = ""
var LogFile = ""
var logger *gotils.Logger
var IP4Hostname = ""
var IP6Hostname = ""
var Port = "8181"
