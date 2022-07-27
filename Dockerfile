# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18 AS build

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build

##
## Deploy
##
FROM debian:stable

WORKDIR /GHost

COPY --from=build /app /GHost

RUN useradd -s /bin/sh -d /GHost GHost
USER GHost

EXPOSE 8000

ENTRYPOINT [ "./GHost" ]

## can be run after build with:
## sudo docker run -p 8000:8000 -v ~/.kube/config:/GHost/.kube/config ghost