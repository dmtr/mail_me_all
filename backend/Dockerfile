FROM golang:1.12.7-alpine3.10 as builder
RUN apk update && apk add git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o mailmeapp ./main

FROM alpine:3.10.1
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /app/mailmeapp /app/
WORKDIR /app
EXPOSE 8080
CMD ["/app/mailmeapp"]
