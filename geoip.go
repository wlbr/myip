package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	geoip "github.com/oschwald/maxminddb-golang"
)

func dialTCP4(network, addr string) (net.Conn, error) {
	return net.Dial("tcp4", addr)
}

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

	pip := net.ParseIP(ip)

	// netip, err := netip.ParseAddr(ip)
	// if err != nil {
	// 	logger.Error("Error parsing ip string: %s", ip, err)
	// } else {
	// 	logger.Info("Parsed ip1: %s", pip)
	// 	logger.Info("Parsed ip2: %s", netip)
	// }
	var rec *record = &record{}

	err = db.Lookup(pip, &rec)
	if err != nil {
		logger.Error("Error finding city for ip %s: %s", pip, err)
		return nil, err
	} else {
		logger.Info("Found city for ip: %s", rec.City.Names)
	}
	return rec, err
}
