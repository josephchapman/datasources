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

## Custom Settings

Check the container image's default `ENTRYPOINT` and `CMD` values:
```bash
docker image inspect otel/opentelemetry-collector-contrib \
| jq '.[].Config | {Entrypoint: .Entrypoint, Cmd: .Cmd}'
```
```
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
┃ ┌─────────────────────────────▼────────┐ ┃ That configuration listends for `otlp` on `:3100`,
┃ │ server:                              │ ┃ and `compose.yaml` exposes the port:
┃ │   http_listen_port: 3100             │ ┃     loki:
┃ └──────────────────────────────────────┘ ┃       ports:
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛         - "3100:3100"
```
