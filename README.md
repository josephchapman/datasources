# datasources

Library of Prometheus Exporters




## Docker compose

Run the example stack:
```bash
cd example-stack/
docker compose up
```

Exporter:
```
http://127.0.0.1:2112/metrics
```



Prometheus:
```
http://localhost:9090
```

Sample query:
```
weather_cloud_cover_percent{location="London"}
```
Sample result:
```
100
```


Grafana
```bash
http://localhost:3000
```