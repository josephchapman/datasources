# datasources

Library of metrics exporters and loggers for use as datasources


## Docker compose

As the containers in this stack run as the non-root user, they will not be able to read files with restrictive umasks (e.g., 077). Before running the stack, ensure that the files and directories are readable to non-root users:
```bash
find ./example-stack -type d -exec chmod 0755 {} \;
find ./example-stack -type f -exec chmod 0644 {} \;
```

Run the example stack:
```bash
cd example-stack/
docker compose up
```

## Utilities

### Prometheus Exporter Weather

```
http://127.0.0.1:2112/metrics
```

### Prometheus
```
http://localhost:9090
```
Test query:
```
weather_temperature_actual_celcius
```

### InfluxDB
```
http://localhost:8086
```
Login:
```
User     : my-user
Password : my-password
```

Query data (Data Explorer):
- Custom Time Range = 5yr
- Script Editor:
```flux
from(bucket: "my-bucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "rating")
  |> filter(fn: (r) => r["player"] == "annacramling")
  |> filter(fn: (r) => r["time_class"] == "blitz")
  |> yield(name: "result")
```

### Loki

Test data:
```bash
$ curl -sX POST \
    http://127.0.0.1:3100/loki/api/v1/push \
    -H 'Content-Type: application/json' \
    --data-raw "
      {
        \"streams\": [
          {
            \"stream\": {
              \"job\": \"testjob\"
            },
            \"values\": [
              [
                \"$(date +%s)000000000\",
                \"it works\"
              ]
            ]
          }
        ]
      }"
```

Confirm test data was received:
```bash
$ curl -s \
    http://127.0.0.1:3100/loki/api/v1/query_range \
    --data-urlencode 'query={job="testjob"}' \
| jq '.data.result'
```
```
[
  {
    "stream": {
      "detected_level": "unknown",
      "job": "testjob",
      "service_name": "testjob"
    },
    "values": [
      [
        "1764277360000000000",
        "it works"
      ]
    ]
  }
]
```

### Grafana
```
http://localhost:3000
```
Login:
```
User     : admin
Password : grafana
```

#### Dashboards

Weather:
![Screenshot](dashboard_weather.png?raw=true "Screenshot of 'Weather' dashboard")

Speedrun:
![Screenshot](dashboard_speedrun.png?raw=true "Screenshot of 'Speedrun' dashboard")

Chess (after running the `replicator-chess` app):
![Screenshot](dashboard_chess.png?raw=true "Screenshot of 'Chess' dashboard")
