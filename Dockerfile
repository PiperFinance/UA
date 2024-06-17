ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
ENV PORT=8000
COPY  ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./src ./src
RUN go build -o ./app ./src/main.go

FROM alpine:latest


RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
COPY  entrypoint.sh .
# COPY ./src/data ./data 

RUN apk update && apk add ca-certificates unzip curl tzdata \
    && cd /tmp \ 
    && curl -OLSs https://github.com/sosedoff/pgweb/releases/download/v0.14.1/pgweb_linux_amd64.zip \
    && unzip pgweb_linux_amd64.zip \
    && mv pgweb_linux_amd64 /api/pgweb \
    && rm -rf /var/bs/log/ | true \ 
    && mkdir -p /var/bs/log/ \ 
    && touch /var/bs/log/err.log \ 
    && touch /var/bs/log/debug.log \
    && rm -rf /var/cache/apk/*


EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
