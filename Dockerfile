FROM golang:1.15.2-alpine3.12 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest AS production
COPY --from=builder /app .
EXPOSE 8086
CMD ["./xnabe-homework"]
