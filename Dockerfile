FROM golang:1.12.7-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /app
ADD . .
RUN apk add git
RUN go build

FROM alpine:3.10.0

RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=builder /app/casabot .
HEALTHCHECK --retries=1 CMD ps aux | grep -q [c]asabot

CMD ["/casabot"]
