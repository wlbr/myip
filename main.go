package main

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	geoip2 "github.com/oschwald/geoip2-golang"
)

//go:generate templify -o myip.go myip.tpl

const GEOIPFILENAME = "GeoLite2-City.mmdb"

const GEOIPURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"

//const GEOIPURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz"
const ETAGFILE = GEOIPFILENAME + ".etag"

type MyIpDate struct {
	Time                     string
	Req                      *http.Request
	IP                       string
	Geo                      *geoip2.City
	City                     string
	Country                  string
	GeoIpFileLastUpdate      string
	GeoIpFileLastUpdateCheck string
}

func NewMyIpDate(r *http.Request, ip string, geo *geoip2.City) *MyIpDate {
	return &MyIpDate{Time: time.Now().Format("January 2, 2006 15:04:05"), Req: r, IP: ip,
		Geo: geo, City: geo.City.Names["en"], Country: geo.Country.Names["en"]}
}

func NewMyIpDateWithUpdate(r *http.Request, ip string, geo *geoip2.City, lastupdate time.Time, lastcheck time.Time) *MyIpDate {
	gipd := NewMyIpDate(r, ip, geo)
	gipd.GeoIpFileLastUpdate = lastupdate.Format("January 2, 2006 15:04:05")
	gipd.GeoIpFileLastUpdateCheck = lastcheck.Format("January 2, 2006 15:04:05")
	return gipd
}

func gentemplate() *template.Template {
	tpl := myipTemplate()
	tmpl, _ := template.New("myip").Parse(tpl)
	return tmpl
}

/*
	checkDownload checks if there is a new GeoIpDatabase available on Maxminds site.
	To save bandwidth it checks the etag on the server (comparing this to a local etag file).
	The db will be downloaded only if the two etags differ.
	The etag is checked on the server only if the last check is more than a day ago.
*/

func checkDownload(url string, file string, c chan bool) {
	var download bool
	var etag string

	etagfildate, err := os.Stat(ETAGFILE)

	if err != nil || time.Now().After(etagfildate.ModTime().AddDate(0, 0, 1)) {
		//Checking Etag, as no etag file found or older than 1 day.
		download = false
		head, err := http.Head(url)
		if err != nil {
			print("Error retrieving GeoIpFile. \n" + err.Error())
		}
		if head.Status != "200 OK" {
			print("Error retrieving GeoIpFile. \nUrl: " + url + "\nStatus: " + head.Status)
			download = false
		} else {
			etag := head.Header.Get("Etag")
			if etag != "" { //etag in header
				fileetag, err := ioutil.ReadFile(ETAGFILE)
				if err != nil { //no old etag found, download
					download = true
				}
				os.Chtimes(ETAGFILE, time.Now(), time.Now())
				if etag != string(fileetag) { //old etag differs from servers one, so download
					download = true
				} else { //old filtype same as servers, do not download
					download = false
				}
			} else { //no etag in header, always download
				download = true
			}
		}
	} else {
		//Not checking Etaga, last check less than 1 day ago
		download = false
	}

	if download {
		fmt.Println("Downloading.")
		out, err := ioutil.TempFile(".", "myip-geoiptmp-")
		defer out.Close()
		if err == nil {
			resp, _ := http.Get(url)
			etag = resp.Header.Get("Etag")
			serverdate, derr := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
			r, _ := gzip.NewReader(resp.Body)
			defer resp.Body.Close()
			if _, err := io.Copy(out, r); err != nil {
				panic("Cannot read from stream!")
			} else {
				os.Rename(out.Name(), file)
				etagout, _ := ioutil.TempFile(".", "myip-etagtmp-")
				defer etagout.Close()
				etagout.Write([]byte(etag))
				os.Rename(etagout.Name(), ETAGFILE)
				if derr == nil {
					os.Chtimes(file, time.Now(), serverdate)
				}
			}
		} else {
			panic("Cannot create file " + file + out.Name() + "\n" + err.Error())
		}

	} else {
		fmt.Println("Not downloading, still the same etag.")
	}
	c <- true
}

func getGeoIp(ip string, w http.ResponseWriter) (*geoip2.City, error) {
	db, err := geoip2.Open(GEOIPFILENAME)

	if err != nil {
		fmt.Fprintf(w, "Error finding GeoIpFile: %s <br> \n\n", err)
		return &geoip2.City{}, err
		//panic("Error open GeoIp: %s\n" + err.Error())
	}
	defer db.Close()

	// If you are using strings that may be invalid, check that ip is not nil
	pip := net.ParseIP(ip)

	record, err := db.City(pip)
	if err != nil {
		fmt.Fprintf(w, "Error finding city: %s\n", err)
		return &geoip2.City{}, err
		//panic("Error find City: %s\n" + err.Error())
	}
	return record, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	record, err := getGeoIp(host, w)
	if err != nil {
		w.Header().Set("Refresh", fmt.Sprintf("5;url=%s", r.RequestURI))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"1;url=%s\"/>", r.RequestURI)
		fmt.Fprintf(w, "Downloading GeoIPData...\n")
	} else {
		geofiledate, err := os.Stat(GEOIPFILENAME)
		etagfiledate, eerr := os.Stat(ETAGFILE)
		if err == nil && eerr == nil {
			err = gentemplate().Execute(w, NewMyIpDateWithUpdate(r, host, record, geofiledate.ModTime(), etagfiledate.ModTime()))
		} else {
			err = gentemplate().Execute(w, NewMyIpDate(r, host, record))
		}
		if err != nil {
			fmt.Fprintf(w, "Error template: %s\n", err)
		}
	}
}

func main() {
	d := make(chan bool, 10)
	go checkDownload(GEOIPURL, GEOIPFILENAME, d)

	runtime.GOMAXPROCS(runtime.NumCPU()) // use all CPU cores
	n := runtime.NumGoroutine() + 1      // initial number of Goroutines

	// install signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	// Spawn request handler
	go func() {
		err := fcgi.Serve(nil, http.HandlerFunc(handler))
		if err != nil {
			fmt.Println("Not in fcgi mode so not spawing handler.")
			c <- syscall.SIGTERM
			//panic(err)
		}
	}()

	// catch signal
	_ = <-d

	// give pending requests in fcgi.Serve() some time to enter the request handler
	time.Sleep(time.Millisecond * 100)

	// wait at most 3 seconds for request handlers to finish
	//inc ase we are downloading the GeoIPFile that may take a while
	for i := 0; i < 30; i++ {
		if runtime.NumGoroutine() <= n {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}

	// catch finished downloader signal
	_ = <-c
}
