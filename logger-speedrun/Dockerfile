FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM golang:1.23.5 AS build
WORKDIR /src
COPY src/* ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/logger-speedrun .

FROM scratch
ENV SPEEDRUN_LEADERBOARDS='[ \
  { \
    "game": "j1l9qz1g", \
    "category": "z275w5k0", \
    "values": { \
      "p854r2vl": "5q85yy6q" \
    } \
  }, \
  { \
    "game": "76rljl68", \
    "category": "jdzmjxkv", \
    "values": { \
      "6nj463vn": "01388y3q" \
    } \
  }, \
  { \
    "game": "o1y9wo6q", \
    "category": "wkpoo02r", \
    "values": { \
      "e8m7em86": "9qj7z0oq" \
    } \
  }, \
  { \
    "game": "o1y9wo6q", \
    "category": "7dgrrxk4", \
    "values": { \
      "e8m7em86": "9qj7z0oq" \
    } \
  } \
]'
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/logger-speedrun /usr/local/bin/logger-speedrun
CMD ["/usr/local/bin/logger-speedrun"]