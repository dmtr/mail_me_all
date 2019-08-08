FROM golang:1.12.7-alpine3.10
RUN apk update && apk add git
RUN mkdir /app
ADD . /app
WORKDIR /app
EXPOSE 8080
RUN go build -o mailmeapp ./main
CMD ["/app/mailmeapp"]
