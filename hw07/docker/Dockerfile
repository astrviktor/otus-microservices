FROM golang:1.19.4-buster as build

ENV BIN_FILE /opt/orderservice
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

LABEL SERVICE="orderservice"
LABEL MAINTAINERS="astrviktor@gmail.com"

ENV BIN_FILE "/opt/orderservice"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}






