# myip
Simple external IP feedback incl. geo location with autoupdate and map display. Runs as fcgi.

See http://wlbr.de/fcgi-bin/myip

##Configuration

Add a configuration 'config.txt" file at build time to Makefiles directory to supply 

   - the site name and analytics id to track the accesses with Google analytics.
   - an upload address to copy the binary to with the deploy make target. Will try an rsync, so take care for the ssh keys.
   - a test url to call with a browser. Simply uses the OS 'open' command usually should start a browser.

Use the '|' as a field separator.

###Example config.txt (not working data):

    analytics|UA-22442552-1|mysite.net
    uploadaddress|myserver:~/fcgi-bin/
    testurl|http://mysite.net/fcgi-bin/myip 
    

