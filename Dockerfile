FROM golang:1.19-alpine3.17 AS build

ADD . ./app
WORKDIR /app
RUN apk update && apk add git curl gcc musl-dev

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN go build -o .

FROM alpine:3.17.0
RUN apk add --update --no-cache ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /app/ethereum-connector /
COPY --from=build /app/scripts/ /scripts/

ENTRYPOINT ./ethereum-connector