<!doctype html>
<html>
<head>
  {{/*<meta http-equiv="Content-Type" content="text/html; charset=us-ascii">
  <link rel="icon" type="image/vnd.microsoft.icon" href="/wlbr/static/res/favicon.ico">*/}}
  <link rel="stylesheet" href="https://cdn.rawgit.com/openlayers/openlayers.github.io/master/en/v5.1.3/css/ol.css" type="text/css">
  {{/*<title>Wolbitest</title>*/}}

  <style>
      .map {
        height: 350px;
        width: 800px;
      }
      body, code { font-family: Arial, Verdana, sans-serif}
  </style>

  <script src="https://cdn.rawgit.com/openlayers/openlayers.github.io/master/en/v5.1.3/build/ol.js"></script>


{{/*if .GoogleAnalyticsId}}
<script type="text/javascript">

  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', '{{.GoogleAnalyticsId}}']);
  _gaq.push(['_setDomainName', '{{.GoogleAnalyticsSite}}']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();

</script>
{{end*/}}

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
      <b>Request IP  : </b>{{.RequestIP}}<br>
      <br>
      <b>Hostnames   : </b>
       {{range $key, $value := .LookupHostnames}}
         {{$value}} &nbsp; &nbsp; &nbsp; &nbsp;
       {{end}}<br>
       <b>IP addresses: </b>
       {{range $key, $value := .LookupIPs}}
         {{$value}} &nbsp; &nbsp; &nbsp; &nbsp;
       {{end}}<br>
      <br>
      <b>City: </b>{{.City}}, {{.Country}} <br>
      <b>Coordinates: </b>{{.Geo.Location.Latitude}}, {{.Geo.Location.Longitude}} <br>
    </p>

    <div id="map" class="map"></div>


   <script type="text/javascript">
      var map = new ol.Map({
        target: 'map',
        loadTilesWhileInteracting: true,
        layers: [
          new ol.layer.Tile({
            source: new ol.source.OSM()
          })
        ],
        view: new ol.View({
          center: ol.proj.fromLonLat([{{.Geo.Location.Longitude}}, {{.Geo.Location.Latitude}}], 'EPSG:3857'),
          zoom: 10
        }),
      });
    </script>

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

  </div>
  <p><br><br><br><br><br><br>This website includes GeoLite2 data created by MaxMind, available from
  <a href="http://www.maxmind.com">http://www.maxmind.com</a>.</p>

</body>
</html>
