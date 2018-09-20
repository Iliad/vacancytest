FROM golang:1.11-alpine as builder
WORKDIR src/github.com/Iliad/vacancytest
COPY . .
RUN go build -o /tmp/vacancytest ./cmd/vacancytest/

FROM alpine:3.7 as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM alpine:3.7
# app
COPY --from=builder /tmp/vacancytest /
# migrations
COPY migrations /migrations
COPY config.yaml /config.yaml
# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/vacancytest"]
