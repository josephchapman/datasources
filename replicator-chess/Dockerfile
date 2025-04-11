FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM golang:1.23.5 AS build
COPY cmn/ /go/src/cmn/
COPY replicator-chess/src/ /go/src/replicator-chess/src/
WORKDIR /go/src/replicator-chess/src/
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/replicator-chess .

FROM scratch
ENV PLAYERS='[ \
  "annacramling", \
  "AlexandraBotez" \
]'
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/replicator-chess /usr/local/bin/replicator-chess
ENTRYPOINT ["/usr/local/bin/replicator-chess"]