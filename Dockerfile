FROM golang:1.25.0-alpine3.22

ARG COMMAND

ENV binary $COMMAND

WORKDIR /go/bin

COPY ${COMMAND} /go/bin

CMD ["sh", "-c", "${binary}"]
