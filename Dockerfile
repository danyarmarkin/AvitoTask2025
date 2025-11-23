FROM golang:1.25.3-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY . .

RUN make build


FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/bin/pr_service /app/pr_service

WORKDIR /app

CMD ["/app/pr_service"]