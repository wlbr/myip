package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getIP(r *http.Request) (string, error) {
	var reqip string
	var e error
	keys, ok := r.URL.Query()["ip"]
	if ok && len(keys) >= 1 {
		reqip = keys[0]
		logger.Info("Got ip %s from request url.", reqip)
	} else {
		reqip = r.Header.Get("X-Forwarded-For")
		if reqip != "" { // ip from header, could be a cdn
			logger.Info("Got ip from X-Forwarded-For-Header")
			cfip := r.Header.Get("Cf-Connecting-Ip")
			if cfip != "" && cfip != reqip {
				logger.Warn("X-Forwarded-For-Header ip and Cf-Connecting-Ip differ. %s != %s", reqip, cfip)
			}
		} else {
			reqip = r.Header.Get("Cf-Connecting-Ip")
			if reqip != "" { //reading cloudflare header
				logger.Info("Got request through cloudflare, connecting ip is: %s", reqip)
			} else { // direct access without cloudflare
				reqip, _, e = net.SplitHostPort(r.RemoteAddr)
				if e != nil {
					logger.Error("Error getting request ip. %s", e)
				}
			}
		}
	}
	return reqip, e
}

func rawSubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Starting rawSubHandler handler.")
	reqip, _ := getIP(r)
	logger.Info("Got raw ip request %s.", reqip)
	// Checking for referer. If the referer is set, then it will be used for CORS.
	// This probably not THAT much better than a wildcard ;-)
	refererurl, e := url.Parse(r.Referer())
	if e != nil || r.Referer() == "" {
		w.WriteHeader(http.StatusLocked)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", refererurl.Scheme+"://"+refererurl.Hostname())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, reqip)
	}
}

func completeProtocol(r *http.Request, site string) string {
	if !(strings.HasPrefix(site, "http://") || strings.HasPrefix(site, "https://")) {
		site = r.URL.Scheme + "://" + site
	}
	return site
}

func completePath(r *http.Request, site string) string {
	site = site + r.URL.Path + "?raw"
	return site
}

func fullTemplateSubHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Starting fullTemplateHandler handler.")

	reqip, _ := getIP(r)
	logger.Info("Got request from IP %s", reqip)

	ip4url := completePath(r, completeProtocol(r, IP4Hostname))
	ip6url := completePath(r, completeProtocol(r, IP6Hostname))
	logger.Info("hostnames:  %s - %s", ip4url, ip6url)

	lookuphosts, err := net.LookupAddr(reqip)
	var lookupips []string
	if err != nil {
		logger.Error("Error resolving ip address: %s error: %s", reqip, err)
	} else {
		for _, host := range lookuphosts {
			logger.Info("Resolving IPs for hostname: %s", host)
			ips, err := net.LookupHost(host)
			if err != nil {
				logger.Error("Error looking up host. %s error: %s", lookuphosts, err)
			} else {
				lookupips = append(lookupips, ips...)
			}
		}
	}

	record, err := getGeoIp(reqip, w)
	if err != nil {
		logger.Error("Error getting GeoIP for %s : %s", reqip, err)
		w.Header().Set("Refresh", fmt.Sprintf("5;url=%s", r.RequestURI))
		fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;url=%s\"/>", r.RequestURI)
		fmt.Fprintf(w, "Downloading GeoIPDatabase ...\n")
	} else {
		geofiledate, err := os.Stat(GEOIPFILENAME)
		etagfiledate, eerr := os.Stat(GEOIPFILENAME + ".etag")
		if err == nil && eerr == nil {
			ipdate := NewMyIpDateWithUpdate(r, reqip, lookupips, lookuphosts, record, AnalyticsID, geofiledate.ModTime(), etagfiledate.ModTime(), ip4url, ip6url)
			err = gentemplate().Execute(w, ipdate)
		} else {
			ipdate := NewMyIpDate(r, reqip, lookupips, lookuphosts, record, AnalyticsID, ip4url, ip6url)
			err = gentemplate().Execute(w, ipdate)
		}
		if err != nil {
			fmt.Fprintf(w, "Error template: %s\n", err)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["raw"]
	if ok && len(keys) >= 1 {
		rawSubHandler(w, r)
	} else {
		fullTemplateSubHandler(w, r)
	}
}
