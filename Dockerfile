FROM golang:alpine AS builder

RUN apk update && apk add --no-cache gcc g++ git ca-certificates pkgconfig tzdata && update-ca-certificates

# Create non-privileged user
RUN adduser -D -g '' appuser
 
COPY . .
RUN go get -d -v

RUN GOOS=linux GOARCH=arm GOARM=7 go build -tags musl -a -installsuffix cgo -ldflags '-extldflags "-static" -w -s' -o /go/bin/binance_feed

# Make image smaller...
FROM scratch

LABEL app=${APP_NAME}
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

# Timezone
ENV TZ=Europe/London

COPY --from=builder /go/bin /go/bin

# Run the binary
CMD ["/go/bin/binance_feed"]
