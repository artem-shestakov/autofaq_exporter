FROM golang:1.22.0-alpine3.18
WORKDIR /app
COPY . .
RUN go build -o autofaq-exporter ./cmd/autofaq-exporter/main.go

FROM alpine:3.18
WORKDIR /autofaq-exporter
COPY --from=0 /app/autofaq-exporter .
CMD [ "./autofaq-exporter" ]
