<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
"http://www.w3.org/TR/html4/loose.dtd">

<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=us-ascii">
  {{/*  <link rel="stylesheet" type="text/css" media="screen" href="//google-code-prettify.googlecode.com/svn/trunk/styles/sunburst.css">
  <style type="text/css" media="screen">@import "/wlbr/static/res/base.css";</style> */}}
  <link rel="icon" type="image/vnd.microsoft.icon" href="/wlbr/static/res/favicon.ico">	  
  <title>Wolbitest</title>
  {{/* <style>
    p { 
        margin-bottom: 2eM;
        margin-bottom: 2eM;
       }
  </style> */}}
  
<style>#gmap_canvas img{max-width:none!important;background:none!important;}
</style>

<script type="text/javascript" src="https://maps.google.com/maps/api/js?sensor=false"></script>

<script type="text/javascript"> 
    function init_map()
      { 
        var myOptions = {zoom:15,center:new google.maps.LatLng({{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}}),mapTypeId: google.maps.MapTypeId.ROADMAP};
        map = new google.maps.Map(document.getElementById("gmap_canvas"), myOptions);
        marker = new google.maps.Marker({map: map,position: new google.maps.LatLng({{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}})});
        infowindow = new google.maps.InfoWindow({content:"<span style='height:auto !important; display:block; white-space:nowrap; overflow:hidden !important;'><strong style='font-weight:400;'>{{.City}}</strong></span>" });
        google.maps.event.addListener(marker, "click", function()    
           {infowindow.open(map,marker);});
        infowindow.open(map,marker);
	   }
     google.maps.event.addDomListener(window, 'load', init_map);
</script>

  
  
</head>

<body>
  <div id="Center">
    <div id="Header">
      <a href="/fcgi-bin/myip" title="MyIP">MyIP</a>
    </div>

    <div id="Content">
    {{.Time}}
       
    <p>
      <b>GeoDatabase updated on: </b>{{.GeoIpFileLastUpdate}}<br>
      <b>GeoDatabase last update-check: </b>{{.GeoIpFileLastUpdateCheck}}<br>
      <br>
      <b>IP: </b>{{.IP}}<br>
      <b>City: </b>{{.City}}, {{.Country}} <br>
      <b>Coordinates: </b>{{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}} <br>
    </p>
    <div style="overflow:hidden;height:250px;width:650x;">
     <div id="gmap_canvas" style="height:250px;width:650px;"></div>
     <a class="google-map-code" href="https://www.map-embed.com" id="get-map-data">https://www.map-embed.com</a>
    </div>
    
	 <p><br>
      <b>Protocol: </b>{{.Req.Proto}}<br>
      <b>Method: </b>{{.Req.Method}}<br>
      <b>UserAgent: </b>{{.Req.UserAgent}}<br>
     </p> 
	
    {{with .Req}}
    <p><br>
      <b>Headers</b><br>
      <code>
       {{range $key, $value := .Header}}
         <b>{{$key}}: </b>{{$value}}<br>
       {{end}}
      </code>      
     </p> 
    {{end}}
    </div>

    <div id="Sidebar">
       {{/* <div id="Menu">
       <a href="http://www.alistapart.com/" title="A List Apart: For People Who Make Websites">A List Apart</a><br>
        <a href="http://www.alistapart.com/" title="A List Apart: For People Who Make Websites">A List Apart</a><br>
      </div> */}}
    </div>
  </div>
  <p><br><br><br><br><br><br>This website includes GeoLite2 data created by MaxMind, available from
<a href="http://www.maxmind.com">http://www.maxmind.com</a>.</p>
</body>
</html>