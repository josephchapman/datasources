FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM golang:1.23.5 AS build
COPY cmn/ /go/src/cmn/
COPY logger-speedrun/src/ /go/src/logger-speedrun/src/
WORKDIR /go/src/logger-speedrun/src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/logger-speedrun .

FROM scratch
ENV SPEEDRUN_LEADERBOARDS='[ \
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
ENTRYPOINT ["/usr/local/bin/logger-speedrun"]