FROM golang:1.19.4-buster as build

ENV BIN_FILE "/opt/auth/auth"
ENV CODE_DIR "/go/src/app"

COPY ./app ${CODE_DIR}
WORKDIR ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} ./cmd/main.go

FROM scratch
#FROM ubuntu:22.04

LABEL SERVICE="auth"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/auth/auth /srv/auth

WORKDIR /srv
ENTRYPOINT ["/srv/auth"]

