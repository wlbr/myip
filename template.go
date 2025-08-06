package main

import (
	"html/template"
	"net/http"
	"time"
)

type MyIpDate struct {
	Time                     string
	Req                      *http.Request
	RequestIP                string
	LookupIPs                []string
	LookupHostnames          []string
	Geo                      *record
	City                     string
	Country                  string
	GeoIpFileLastUpdate      string
	GeoIpFileLastUpdateCheck string
	GoogleAnalyticsId        string
}

type record struct {
	Country struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Location struct {
		AccuracyRadius int     `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

func NewMyIpDate(r *http.Request, ip string, lookupips, lookuphostnames []string, geo *record, analyticsid string) *MyIpDate {
	return &MyIpDate{Time: time.Now().Format("January 2, 2006 15:04:05"),
		Req: r, RequestIP: ip, LookupIPs: lookupips, LookupHostnames: lookuphostnames,
		Geo: geo, City: geo.City.Names["en"], Country: geo.Country.Names["en"],
		GoogleAnalyticsId: analyticsid}
}

func NewMyIpDateWithUpdate(r *http.Request, ip string, lookupips, lookuphostnames []string, geo *record, analyticsid string, lastupdate time.Time, lastcheck time.Time) *MyIpDate {
	gipd := NewMyIpDate(r, ip, lookupips, lookuphostnames, geo, analyticsid)
	gipd.GeoIpFileLastUpdate = lastupdate.Format("January 2, 2006 15:04:05")
	gipd.GeoIpFileLastUpdateCheck = lastcheck.Format("January 2, 2006 15:04:05")
	return gipd
}

func gentemplate() *template.Template {
	tpl := myipTemplate()
	tmpl, _ := template.New("myip").Parse(tpl)
	return tmpl
}
