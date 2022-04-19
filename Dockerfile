FROM golang:1.17-buster as build

WORKDIR /base
ADD . /base/
RUN go build -o /base/app .


LABEL io.airbyte.version=0.0.1
LABEL io.airbyte.name=airbyte/source

ENTRYPOINT ["/base/app"]