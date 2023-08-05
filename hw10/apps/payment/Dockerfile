FROM golang:1.19.4-buster as build

ENV BIN_FILE "/opt/payment/payment"
ENV CODE_DIR "/go/src/app"

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o ${BIN_FILE} main.go

FROM scratch

LABEL SERVICE="payment"
LABEL MAINTAINERS="astrviktor@gmail.com"

COPY --from=build /opt/payment/payment /srv/payment

WORKDIR /srv
ENTRYPOINT ["/srv/payment"]





