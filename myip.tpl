<!doctype html>
<html>
<head>
  {{/*<meta http-equiv="Content-Type" content="text/html; charset=us-ascii">
  <link rel="icon" type="image/vnd.microsoft.icon" href="/wlbr/static/res/favicon.ico">*/}}
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/openlayers/openlayers.github.io@main/dist/en/v7.1.0/ol/ol.css" />
  <script src="https://cdn.jsdelivr.net/gh/openlayers/openlayers.github.io@main/dist/en/v7.1.0/ol/dist/ol.js"></script>
  {{/*<title>Wolbitest</title>*/}}

  <style>
      .map {
        height: 350px;
        width: 800px;
      }
      body, code { font-family: Arial, Verdana, sans-serif}
  </style>


{{if .GoogleAnalyticsId}}
<!-- Google tag (gtag.js) -->
<script async src="https://www.googletagmanager.com/gtag/js?id={{.GoogleAnalyticsId}}"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', '{{.GoogleAnalyticsId}}');
</script>
{{end}}



</head>

<body>

  <div id="Center">
    <div id="Header">
      <a href="." title="MyIP">MyIP</a>
    </div>

    <div id="Content">
    {{.Time}}

    <p>
      <b>GeoDatabase updated on: </b>{{.GeoIpFileLastUpdate}}<br>
      <b>GeoDatabase last update-check: </b>{{.GeoIpFileLastUpdateCheck}}<br>
      <br>
      <b>Request IP  : </b>{{.RequestIP}}<br>
      </p>
    <p>
     <b>Your Hostnames   : </b>
       {{range $key, $value:= .LookupHostnames}}
         {{$value}} &nbsp; &nbsp; &nbsp; &nbsp;
       {{end}}<br>
       <b>Your hostnames IP addresses: </b>
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
