FROM alpine:latest

RUN apk add --no-cache git go

WORKDIR /app

COPY . /app

RUN go build -ldflags="-s -w"

CMD ["./app"]
