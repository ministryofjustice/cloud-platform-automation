FROM golang:1.23.6-alpine3.21

ARG COMMAND

ENV binary $COMMAND

WORKDIR /go/bin

COPY ${COMMAND} /go/bin

CMD ["sh", "-c", "${binary}"]