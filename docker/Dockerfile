FROM golang:1.13-alpine AS build-env
ADD . /src
WORKDIR /src/cmd/scan_mail
RUN apk add git gcc musl-dev && go build

FROM alpine:latest

COPY --from=build-env /src/cmd/scan_mail/scan_mail /opt/scan_mail
COPY docker/entrypoint.sh /entrypoint.sh

RUN apk add postfix cyrus-sasl cyrus-sasl-plain cyrus-sasl-login cyrus-sasl-digestmd5 cyrus-sasl-ntlm cyrus-sasl-crammd5 rsyslog supervisor

ENTRYPOINT 
CMD ["/entrypoint.sh"]
