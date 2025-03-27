# datasources

Library of metrics exporters and loggers for use as datasources


## Docker compose

Run the example stack:
```bash
cd example-stack/
docker compose up
```

## Utilities

### Prometheus
```
http://localhost:9090
```

### InfluxDB
```bash
http://localhost:8086
```

### Grafana
```bash
http://localhost:3000
```

#### Dashboards

Weather:
![Screenshot](dashboard_weather.png?raw=true "Screenshot of 'Weather' dashboard")

Speedrun:
![Screenshot](dashboard_speedrun.png?raw=true "Screenshot of 'Speedrun' dashboard")

Chess (after running the `replicator-chess` app):
![Screenshot](dashboard_chess.png?raw=true "Screenshot of 'Chess' dashboard")
