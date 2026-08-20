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

#### Writing Metrics

```bash
curl "http://localhost:8086/api/v2/write" \
  --url-query "org=my-org" \
  --url-query "bucket=my-bucket" \
  --url-query "precision=s" \
  -X POST \
  -H "Authorization: Token my-super-secret-auth-token" \
  --data-binary "
    utilities,source=bill  electricity_import=1234,electricity_export=2345,electricity_generation=3456,electricity_consumption=4567,water_import=5678 $(date -d '2026-08-02 12:00 CEST' +%s)
    utilities,source=meter electricity_import=1334,electricity_export=2445,electricity_generation=3556,electricity_consumption=4667,water_import=5778 $(date -d '2026-08-09 12:00 CEST' +%s)
    utilities,source=meter electricity_import=1434,electricity_export=2545,electricity_generation=3656,electricity_consumption=4767,water_import=5878 $(date -d '2026-08-16 12:00 CEST' +%s)
  "
```

#### Clearing Metrics

```bash
curl "http://localhost:8086/api/v2/delete" \
  --url-query "org=my-org" \
  --url-query "bucket=my-bucket" \
  -X POST \
  -H "Authorization: Token my-super-secret-auth-token" \
  -d '{
    "start": "1970-01-01T00:00:00Z",
    "stop": "2070-01-01T00:00:00Z",
    "predicate": "_measurement=\"utilities\""
  }'
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
```json
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

### Interface

```
http://localhost:8501
```

## Custom Settings

Check the container image's default `ENTRYPOINT` and `CMD` values:
```bash
docker image inspect otel/opentelemetry-collector-contrib \
| jq '.[].Config | {Entrypoint: .Entrypoint, Cmd: .Cmd}'
```
```json
{
  "Entrypoint": [
    "/otelcol-contrib"
  ],
  "Cmd": [
    "--config",
    "/etc/otelcol-contrib/config.yaml"
  ]
}
```

## example-stack Architecture

### Logging

The OpenTelemetry Collector is used to receive, process and export logs.
Docker does not have an OpenTelemetry logging driver, so the Fluent logging driver is used.
OpenTelemetry is configured to listen for Fluent connections.

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ `compose.yaml` defines a YAML Anchor (`&logging`),
┃ Docker Logging Driver                    ┃ which is inherited by services by adding:
┃                                          ┃     logging: *logging
┃ ┌──────────────────────────────────────┐ ┃
┃ │ driver: "fluentd"                    │ ┃ It replaces Docker's default `json-file` with `fluentd`,
┃ │ options:                             │ ┃ and sends entries to a central receiver at:
┃ │   fluentd-address: "localhost:24224" │ ┃     localhost:24224
┃ └─────────────────────────────┬────────┘ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┛
                                │
  [ protocol : fluent_forward ] │
                                │
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┓ `compose.yaml` defines an otel-collector service configured via:
┃ OpenTelemetry Collector       │          ┃     ./otelcol-contib/config.yaml
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ That configuration listens for `fluent_forward` on `:24224`,
┃ │ receiver:                            │ ┃ and `compose.yaml` exposes the port:
┃ │   fluent_forward:                    │ ┃     otel-collector:
┃ │     endpoint: 0.0.0.0:24224          │ ┃       ports:
┃ └─────────────────────────────┬────────┘ ┃         - "24224:24224"
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ It then sends the entries via `otlp_http` to a log store at:
┃ │ exporter:                            │ ┃     loki:3100/otlp
┃ │   otlp_http:                         │ ┃
┃ │     endpoint: http://loki:3100/otlp  │ ┃
┃ └─────────────────────────────┬────────┘ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┛
                                │
       [ protocol : otlp_http ] │
                                │
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┓ `compose.yaml` defines a loki service confirgured via:
┃ Loki                          │          ┃     ./loki/local-config.yaml
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ That configuration listens for `otlp` on `:3100`,
┃ │ server:                              │ ┃ and `compose.yaml` exposes the port:
┃ │   http_listen_port: 3100             │ ┃     loki:
┃ └──────────────────────────────────────┘ ┃       ports:
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛         - "3100:3100"
```

### Metrics

Scrape method:
```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ Applications may expose a `/metrics` endpoint on a port, such as:
┃ Application                              ┃     ports:
┃                                          ┃       - 2112:2112
┃ ┌──────────────────────────────────────┐ ┃
┃ │ :2112/metrics                        │ ┃
┃ └─────────────────────────────┬────────┘ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┛
                                │
                                │
                                │
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┓ `compose.yaml` defines an otel-collector service configured via:
┃ OpenTelemetry Collector       │          ┃     ./otelcol-contib/config.yaml
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ That configuration scrapes metrics.
┃ │ receiver:                            │ ┃
┃ │   prometheus:                        │ ┃
┃ │     config:                          │ ┃
┃ │       scrape_configs:                │ ┃
┃ │         - job_name: weather          │ ┃
┃ │           static_configs:            │ ┃
┃ │             - targets:               │ ┃
┃ │                 - weather:2112       │ ┃
┃ └─────────────────────────────┬────────┘ ┃
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ It then sends the metrics via `prometheusremotewrite` to the metrics store at:
┃ │ exporter:                            │ ┃     prometheus:9090/api/v1/write
┃ │   prometheusremotewrite:             │ ┃
┃ │     endpoint: prometheus/api/v1/write│ ┃
┃ └─────────────────────────────┬────────┘ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┛
                                │
                                │
                                │
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━│━━━━━━━━━━┓ `compose.yaml` defines a Prometheus service confirgured via:
┃ Prometheus                    │          ┃     ./prometheus/prometheus.yaml
┃                               │          ┃
┃ ┌─────────────────────────────▼────────┐ ┃ That configuration listens on:
┃ │ command:                             │ ┃     prometheus:9090/api/v1/write
┃ │   --web.enable-remote-write-receiver │ ┃
┃ └──────────────────────────────────────┘ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```