FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM golang:1.23.5 AS build
WORKDIR /src
COPY src/* ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/prometheus-exporter-weather .

FROM scratch
ENV WEATHER_LOCATIONS='[ \
  { \
    "name": "New York", \
    "latitude": 40.79, \
    "longitude": -73.96, \
    "tzdata": "America/New_York" \
  }, \
  { \
    "name": "London", \
    "latitude": 51.50, \
    "longitude": -0.12, \
    "tzdata": "Europe/London" \
  } \
]'
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/prometheus-exporter-weather /usr/local/bin/prometheus-exporter-weather
CMD ["/usr/local/bin/prometheus-exporter-weather"]