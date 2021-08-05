FROM golang:1.16.4-alpine3.13 as builder
WORKDIR /app
RUN apk update && apk upgrade  -U -a && \
    apk add bash git openssh gcc libc-dev

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
# Build the Go app
RUN go build -o /app/server /app


FROM golang:1.16.4-alpine3.13
RUN apk add --update ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    echo "Asia/Ho_Chi_Minh" > /etc/timezone && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/server /app/server

WORKDIR /app
CMD /app/server
