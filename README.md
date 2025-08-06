# myip

Simple external IP feedback incl. geo location with autoupdate and map display. Runs as fcgi or standalone http server.

See <http://myip.wlbr.de:8282>

## Configuration

Add a 'config.txt" file for build time configuration to the directory of the Makefile to add

- the download URL provided by Maxmind. You will need to register and provide the token given by them
- the analytics id to track the accesses with Google analytics.
- an upload address to copy the binary to when running the deploy make target. Will try an rsync, so take care for the ssh keys.
- the port the server listens to if not in fcgi mode.
- loglevel, that can be one of (Off, Fatal, Error, Warn, Debsug, Info, All).
- logfile, a  be a filename or STDOUT.

Use the '|' as a field separator.

### Example config.txt (not working data)

    downloadurl|https://yourid:yourkey@download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz
    analytics|G-L8JUH7WYB
    loglevel|All
    logfile|myip.log
    port|5150
    uploadaddress|myserver:~/bin/

