/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package main

// myipTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func myipTemplate() string {
	var tmpl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\"\n" +
		"\"http://www.w3.org/TR/html4/loose.dtd\">\n" +
		"\n" +
		"<html>\n" +
		"<head>\n" +
		"  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=us-ascii\">\n" +
		"  {{/*  <link rel=\"stylesheet\" type=\"text/css\" media=\"screen\" href=\"//google-code-prettify.googlecode.com/svn/trunk/styles/sunburst.css\">\n" +
		"  <style type=\"text/css\" media=\"screen\">@import \"/wlbr/static/res/base.css\";</style> */}}\n" +
		"  <link rel=\"icon\" type=\"image/vnd.microsoft.icon\" href=\"/wlbr/static/res/favicon.ico\">\t  \n" +
		"  <title>Wolbitest</title>\n" +
		"  {{/* <style>\n" +
		"    p { \n" +
		"        margin-bottom: 2eM;\n" +
		"        margin-bottom: 2eM;\n" +
		"       }\n" +
		"  </style> */}}\n" +
		"  \n" +
		"<style>#gmap_canvas img{max-width:none!important;background:none!important;}\n" +
		"</style>\n" +
		"\n" +
		"<script type=\"text/javascript\" src=\"https://maps.google.com/maps/api/js?sensor=false\"></script>\n" +
		"\n" +
		"<script type=\"text/javascript\"> \n" +
		"    function init_map()\n" +
		"      { \n" +
		"        var myOptions = {zoom:15,center:new google.maps.LatLng({{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}}),mapTypeId: google.maps.MapTypeId.ROADMAP};\n" +
		"        map = new google.maps.Map(document.getElementById(\"gmap_canvas\"), myOptions);\n" +
		"        marker = new google.maps.Marker({map: map,position: new google.maps.LatLng({{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}})});\n" +
		"        infowindow = new google.maps.InfoWindow({content:\"<span style='height:auto !important; display:block; white-space:nowrap; overflow:hidden !important;'><strong style='font-weight:400;'>{{.City}}</strong></span>\" });\n" +
		"        google.maps.event.addListener(marker, \"click\", function()    \n" +
		"           {infowindow.open(map,marker);});\n" +
		"        infowindow.open(map,marker);\n" +
		"\t   }\n" +
		"     google.maps.event.addDomListener(window, 'load', init_map);\n" +
		"</script>\n" +
		"\n" +
		"  \n" +
		"  \n" +
		"</head>\n" +
		"\n" +
		"<body>\n" +
		"  <div id=\"Center\">\n" +
		"    <div id=\"Header\">\n" +
		"      <a href=\"/fcgi-bin/myip\" title=\"MyIP\">MyIP</a>\n" +
		"    </div>\n" +
		"\n" +
		"    <div id=\"Content\">\n" +
		"    {{.Time}}\n" +
		"       \n" +
		"    <p>\n" +
		"      <b>GeoDatabase updated on: </b>{{.GeoIpFileLastUpdate}}<br>\n" +
		"      <b>GeoDatabase last update-check: </b>{{.GeoIpFileLastUpdateCheck}}<br>\n" +
		"      <br>\n" +
		"      <b>IP: </b>{{.IP}}<br>\n" +
		"      <b>City: </b>{{.City}}, {{.Country}} <br>\n" +
		"      <b>Coordinates: </b>{{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}} <br>\n" +
		"    </p>\n" +
		"    <div style=\"overflow:hidden;height:250px;width:650x;\">\n" +
		"     <div id=\"gmap_canvas\" style=\"height:250px;width:650px;\"></div>\n" +
		"     <a class=\"google-map-code\" href=\"https://www.map-embed.com\" id=\"get-map-data\">https://www.map-embed.com</a>\n" +
		"    </div>\n" +
		"    \n" +
		"\t <p><br>\n" +
		"      <b>Protocol: </b>{{.Req.Proto}}<br>\n" +
		"      <b>Method: </b>{{.Req.Method}}<br>\n" +
		"      <b>UserAgent: </b>{{.Req.UserAgent}}<br>\n" +
		"     </p> \n" +
		"\t\n" +
		"    {{with .Req}}\n" +
		"    <p><br>\n" +
		"      <b>Headers</b><br>\n" +
		"      <code>\n" +
		"       {{range $key, $value := .Header}}\n" +
		"         <b>{{$key}}: </b>{{$value}}<br>\n" +
		"       {{end}}\n" +
		"      </code>      \n" +
		"     </p> \n" +
		"    {{end}}\n" +
		"    </div>\n" +
		"\n" +
		"    <div id=\"Sidebar\">\n" +
		"       {{/* <div id=\"Menu\">\n" +
		"       <a href=\"http://www.alistapart.com/\" title=\"A List Apart: For People Who Make Websites\">A List Apart</a><br>\n" +
		"        <a href=\"http://www.alistapart.com/\" title=\"A List Apart: For People Who Make Websites\">A List Apart</a><br>\n" +
		"      </div> */}}\n" +
		"    </div>\n" +
		"  </div>\n" +
		"  <p><br><br><br><br><br><br>This website includes GeoLite2 data created by MaxMind, available from\n" +
		"<a href=\"http://www.maxmind.com\">http://www.maxmind.com</a>.</p>\n" +
		"</body>\n" +
		"</html>"
	return tmpl
}
