FROM golang:1.12.7-alpine3.10 as builder_base
RUN apk update && apk add git
RUN mkdir /app
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod download

FROM builder_base as builder
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o mailmeapp ./main

FROM alpine:3.10.1
RUN apk add --no-cache tzdata
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /app/mailmeapp /app/
WORKDIR /app
EXPOSE 8080
CMD ["/app/mailmeapp"]
