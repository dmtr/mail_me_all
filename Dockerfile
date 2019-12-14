FROM golang:1.12.7-alpine3.10 as builder_base
RUN apk update && apk add git
RUN mkdir -p /app/backend
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod download

FROM builder_base as builder
COPY ./backend /app/backend
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o mailmeapp ./backend/main

FROM alpine:3.10.1 as service
ARG APP_PORT
RUN apk add --no-cache tzdata ca-certificates \
    && update-ca-certificates 2>/dev/null || true
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /app/mailmeapp /app/
COPY --from=builder /app/backend/templates /app/templates
COPY --from=builder /app/backend/cert/service.pem /app/service.pem	
COPY --from=builder /app/backend/cert/service.key /app/service.key	
WORKDIR /app
EXPOSE ${APP_PORT:-8080}
CMD ["/app/mailmeapp"]

FROM alpine:3.10.1 as crontask
RUN apk add --no-cache tzdata ca-certificates \
    && update-ca-certificates 2>/dev/null || true
COPY --from=builder /app/mailmeapp /app/
COPY --from=builder /app/backend/crontasks/entrypoint.sh /app/entrypoint.sh
COPY --from=builder /app/backend/cert/service.pem /app/service.pem	
COPY --from=builder /app/backend/cert/service.key /app/service.key	
COPY --from=builder /app/backend/templates /app/templates
RUN chmod +x /app/entrypoint.sh
WORKDIR /app
ENTRYPOINT ["sh", "/app/entrypoint.sh"]
CMD [""]
