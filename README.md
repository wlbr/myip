# myip

Simple external IP feedback incl. geo location with autoupdate and map display. Runs as fcgi.

See <http://wlbr.de/fcgi-bin/myip>

## Configuration

Add a 'config.txt" file for build time configuration to the directory of the Makefile to add

- the site name and analytics id to track the accesses with Google analytics.
- an upload address to copy the binary to when running the deploy make target. Will try an rsync, so take care for the ssh keys.
- a test url to open with a browser. Simply uses the OS 'open' command.
- loglevel, that can be one of (Off, Fatal, Error, Warn, Debsug, Info, All).
- logfile, a  be a filename or STDOUT.

Use the '|' as a field separator.

### Example config.txt (not working data)

    analytics|UA-22442552-1|mysite.net
    uploadaddress|myserver:~/fcgi-bin/
    testurl|http://mysite.net/fcgi-bin/myip
    loglevel|All
    logfile|myip.log
    ip4hostname|www4.fcgi.net
    ip6hostname|www6.fcgi.net

