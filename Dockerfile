FROM golang:1.20.5 AS builder
WORKDIR /app
COPY . .
RUN GOOS-linux GOARCH-amd64 go build -o time-tz

FROM scratch
COPY --from=builder /app/time-tz /opt/time-tz
CMD ["/opt/time-tz"]