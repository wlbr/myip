# myip
Simple external IP feedback incl. geo location with autoupdate and map display. Runs as fcgi.

See http://wlbr.de/fcgi-bin/myip

Add a configuration file to add 

   - the site name and analytics id to track the accesses with Google analytics.
   - an upload address to copy the binary with the deploy make target. Will try an rsync, so take care for the ssh keys.
   - a test url to call with an browser. Simply uses the OS 'open' command usually should start a browser.


Example (not working data):

    analytics|UA-22442529-1|mysite.net
    uploadaddress|myserver:~/fcgi-bin/
    testurl|http://mysite.net/fcgi-bin/myip 
    

