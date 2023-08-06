#
## Build
##

FROM golang:alpine AS build

WORKDIR /app

COPY . /app/

RUN go build -o carlog main.go

##
## Deploy
##

FROM alpine as carlog

WORKDIR /

COPY --from=build /app/carlog /carlog
CMD ["/carlog"]
