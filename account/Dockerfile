FROM golang:1.20-alpine AS builder
WORKDIR /source
COPY . /source
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o account-service ./cmd/.

FROM alpine:3.15
RUN mkdir /app
WORKDIR /app
COPY --from=builder /source/account-service /usr/local/bin
COPY --from=builder /source/configs /app/configs
COPY --from=builder /source/migrations /app/migrations
RUN chmod a+x /usr/local/bin/account-service
ENTRYPOINT [ "account-service" ]
