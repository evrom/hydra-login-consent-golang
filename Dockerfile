
FROM golang:1.17-alpine3.15 AS builder

WORKDIR /go/src/github.com/evrom/hydra-login-consent-golang

ADD go.mod go.mod
ADD go.sum go.sum

RUN go mod download

ADD . .

RUN go build -o /usr/bin/hydra-login-consent-golang

FROM alpine:3.15

RUN adduser -S ory -D -u 10000 -s /bin/nologin

COPY --from=builder /usr/bin/hydra-login-consent-golang /usr/bin/hydra-login-consent-golang
ADD . .
EXPOSE 3000

USER ory

ENTRYPOINT ["hydra-login-consent-golang"]
CMD ["serve"]
