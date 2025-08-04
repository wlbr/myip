package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	//geoip2 "github.com/oschwald/geoip2-golang"
	geoip "github.com/oschwald/maxminddb-golang"

	"github.com/wlbr/myip/gotils"
)

//go:generate templify -o myip.go myip.tpl

// GEOIPFILENAME is the local filename to the GeoIp database.
const GEOIPFILENAME = "GeoLite2-City.mmdb"

// GeoIpUrl is the url to download the GeoIp database from. Set by linker.
// var GeoIpUrl = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
var GeoIpUrl = "https://ipinfo.io/data/ipinfo_lite.mmdb?token=YOURTOKEN"

// ETAGFILE is the file of the last downloads etag hash. It will be compared to the
// servers etag to prevent unnessecary downloads.
//const ETAGFILE = GEOIPFILENAME + ".etag"

// AnalyticsId is the ID to be used by Google Analytics. It is set by linker flags.
// If AnalyticsId not set, then the tracking code will be omitted.
var AnalyticsID = ""

// AnalyticsSite is the site name to be used by Google Analytics.
// It is set by linker flags. If AnalyticsSite is not set, then the tracking code will be omitted.
var AnalyticsSite = ""

// If LogLevel not set, then the logging is shut off. Set by linker.
var LogLevel = ""

// If LogFile not set, then the logging is set to os.stdout. Set by linker.
var LogFile = ""
var logger *gotils.Logger

// If LogFile not set, then the logging is set to os.stdout. Set by linker.
var IP4Hostname = ""
var IP6Hostname = ""

// Port is the port to listen on. If not set, then the default is 8181. Set by linker.
var Port = "8181"

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
	GoogleAnalyticsSite      string
	IP4Hostname              string
	IP6Hostname              string
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

func NewMyIpDate(r *http.Request, ip string, lookupips, lookuphostnames []string, geo *record, analyticsid, analyticssite, ip4hostname, ip6hostname string) *MyIpDate {
	return &MyIpDate{Time: time.Now().Format("January 2, 2006 15:04:05"),
		Req: r, RequestIP: ip, LookupIPs: lookupips, LookupHostnames: lookuphostnames,
		Geo: geo, City: geo.City.Names["en"], Country: geo.Country.Names["en"],
		GoogleAnalyticsId: analyticsid, GoogleAnalyticsSite: analyticssite,
		IP4Hostname: ip4hostname, IP6Hostname: ip6hostname}
}

func NewMyIpDateWithUpdate(r *http.Request, ip string, lookupips, lookuphostnames []string, geo *record,
	analyticsid, analyticssite string, lastupdate time.Time, lastcheck time.Time, ip4id, ip6id string) *MyIpDate {
	gipd := NewMyIpDate(r, ip, lookupips, lookuphostnames, geo, analyticsid, analyticssite, ip4id, ip6id)
	gipd.GeoIpFileLastUpdate = lastupdate.Format("January 2, 2006 15:04:05")
	gipd.GeoIpFileLastUpdateCheck = lastcheck.Format("January 2, 2006 15:04:05")
	return gipd
}

func gentemplate() *template.Template {
	tpl := myipTemplate()
	tmpl, _ := template.New("myip").Parse(tpl)
	return tmpl
}

func dialTCP4(network, addr string) (net.Conn, error) {
	return net.Dial("tcp4", addr)
}

/*
	checkDownload checks if there is a new GeoIpDatabase available on Maxminds site.
	To save bandwidth it checks the etag on the server (comparing this to a local etag file).
	The db will be downloaded only if the two etags differ.
	The etag is checked on the server only if the last check is more than a day ago.
*/

func checkDownload(uri string, file string, c chan bool) {
	var download bool
	var etag string
	var etagfile = file + ".etag"

	tr := &http.Transport{Dial: dialTCP4}
	httpclient := &http.Client{Transport: tr}

	etagfildate, err := os.Stat(etagfile)
	if err != nil {
		logger.Info("Error reading etag file %s. %s", etagfile, err.Error())
	}
	if err != nil || time.Now().After(etagfildate.ModTime().AddDate(0, 0, 1)) {
		//Checking Etag, as no etag file found or older than 1 day.
		logger.Info("Checking download from %s", GeoIpUrl)
		download = false
		head, err := httpclient.Head(uri)

		if err != nil {
			logger.Error("Error retrieving GeoIpFile. %s", err.Error())
		}
		if head.Status != "200 OK" {
			logger.Error("Error retrieving GeoIpFile. Url: %s Status: %s", uri, head.Status)
			download = false
		} else {
			etag := head.Header.Get("Etag")
			if etag != "" { //etag in header
				cwd, err := os.Getwd()
				logger.Info("cwd= %s", cwd)
				fileetag, err := os.ReadFile(etagfile)
				if err != nil { //no old etag found, download
					download = true
					logger.Info("No local etag found, downloading. err=%v", err)
				}
				os.Chtimes(etagfile, time.Now(), time.Now())
				if etag != string(fileetag) { //old etag differs from servers one, so download
					download = true
					logger.Info("Etag differs from server, downloading.")
				} else { //old filetype same as servers, do not download
					download = false
					logger.Info("Same etag, not downloading.")
				}
			} else { //no etag in header, always download
				download = true
				logger.Info("No etag found on server, downloading.")
			}
		}
	} else {
		//Not checking Etag, last check less than 1 day ago
		download = false
		logger.Info("Not checking for download, last check less than a day ago.")
	}

	if download {
		logger.Info("Downloading GeoIp database from %s.", GeoIpUrl)

		resp, _ := httpclient.Get(uri)
		etag = resp.Header.Get("Etag")
		serverdate, derr := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
		var r io.ReadCloser
		r, _ = gzip.NewReader(resp.Body)
		defer r.Close()

		// Extract tar.gz and find the .mmdb file
		tr := tar.NewReader(r)
		var mmdbFile *os.File
		//var mmdbFileName string

		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Fatal("Error reading tar: %s", err.Error())
			}

			// Look for .mmdb file
			if strings.HasSuffix(header.Name, ".mmdb") {
				mmdbFile, err = os.CreateTemp(".", "tmp-myip-mmdb-")
				if err != nil {
					logger.Error("Cannot create temp mmdb file for download: %s", err.Error())
				} else {
					defer mmdbFile.Close()
					if _, err := io.Copy(mmdbFile, tr); err != nil {
						logger.Fatal("Cannot extract mmdb file '%s': %s", mmdbFile.Name(), err.Error())
					}
					break
				}
			}
		}
		if mmdbFile == nil {
			logger.Fatal("No .mmdb file found in archive")
		}

		// Move the extracted mmdb file to the target location
		os.Rename(mmdbFile.Name(), file)
		etagout, _ := os.CreateTemp(".", "tmp-myip-mmdb-etag-")
		defer etagout.Close()
		etagout.Write([]byte(etag))
		os.Rename(etagout.Name(), etagfile)
		if derr == nil {
			os.Chtimes(file, time.Now(), serverdate)
		}
	}
	c <- true
}

func getGeoIp(ip string, w http.ResponseWriter) (*record, error) {
	logger.Info("Opening GeoIP database: %s", GEOIPFILENAME)
	db, err := geoip.Open(GEOIPFILENAME)

	if err != nil {
		logger.Error("Error finding GeoIpFile '%s': %s", GEOIPFILENAME, err)
		return nil, err
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that ip is not nil
	//pip := net.ParseIP("93.202.185.104")
	pip := net.ParseIP(ip)

	netip, err := netip.ParseAddr(ip)
	if err != nil {
		logger.Error("Error parsing ip string: %s", ip, err)
	} else {
		logger.Info("Parsed ip1: %s", pip)
		logger.Info("Parsed ip2: %s", netip)
	}
	var record *record = &record{}

	err = db.Lookup(pip, &record)
	if err != nil {
		logger.Error("Error finding city for ip %s: %s", pip, err)
		return nil, err
	} else {
		logger.Info("Found city for ip: %s", record.City.Names)
	}
	return record, err
}

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
		//w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;url=%s\"/>", r.RequestURI)
		fmt.Fprintf(w, "Downloading GeoIPDatabase ...\n")
	} else {
		geofiledate, err := os.Stat(GEOIPFILENAME)
		etagfiledate, eerr := os.Stat(GEOIPFILENAME + ".etag")
		if err == nil && eerr == nil {
			err = gentemplate().Execute(w, NewMyIpDateWithUpdate(r, reqip, lookupips, lookuphosts, record, AnalyticsID,
				AnalyticsSite, geofiledate.ModTime(), etagfiledate.ModTime(), ip4url, ip6url))
		} else {
			err = gentemplate().Execute(w, NewMyIpDate(r, reqip, lookupips, lookuphosts, record, AnalyticsID, AnalyticsSite, ip4url, ip6url))
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

func main() {
	mode := "http"
	llevel, _ := gotils.LogLevelString(LogLevel)
	logger = gotils.NewLogger(LogFile, llevel)
	logger.Info("Starting up")

	d := make(chan bool, 10)
	go checkDownload(GeoIpUrl, GEOIPFILENAME, d)

	runtime.GOMAXPROCS(runtime.NumCPU()) // use all CPU cores
	n := runtime.NumGoroutine() + 1      // initial number of Goroutines

	// install signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	// Spawn request handler

	// go func() {
	// 	err := fcgi.Serve(nil, http.HandlerFunc(handler))
	// 	if err != nil {//
	// 		c <- syscall.SIGTERM
	// 		mode = "http"
	// 		//panic(err)

	// 	} else {
	// 		logger.Info("Spawing fcgi handler.")
	// 		mode = "fcgi"
	// 	}
	// }()

	if mode == "fcgi" {
		// catch signal
		_ = <-d

		// give pending requests in fcgi.Serve() some time to enter the request handler
		time.Sleep(time.Millisecond * 100)

		// wait at most 3 seconds for request handlers to finish
		// incase we are downloading the GeoIPFile that may take a while
		for i := 0; i < 30; i++ {
			if runtime.NumGoroutine() <= n {
				return
			}
			time.Sleep(time.Millisecond * 100)
		}

		// catch finished downloader signal
		_ = <-c
	} else {
		logger.Info("Not in fcgi mode so not spawing handler. Will start http server on port %s.", Port)
		http.HandleFunc("/", handler)
		errsrv := http.ListenAndServe(":"+Port, nil)
		if errsrv != nil {
			logger.Fatal("Error starting http server on port %s: %s", Port, errsrv.Error())
		}
	}
}
