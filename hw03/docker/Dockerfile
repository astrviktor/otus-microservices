FROM golang:1.19.4-buster as build

ENV BIN_FILE /opt/crudservice
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/*

# ------------------------------------
FROM alpine:3.17.0

LABEL SERVICE="crudservice"
LABEL MAINTAINERS="astrviktor@gmail.com"

ENV BIN_FILE "/opt/crudservice"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/crudservice/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}






