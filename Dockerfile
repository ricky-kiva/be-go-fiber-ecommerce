FROM golang:1.21-alpine3.18 AS builder
RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /go/src

ENV GOPROXY=https://goproxy.io,direct
COPY . .

RUN go mod tidy -x
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o main ./

FROM alpine:3.18 AS runner
RUN adduser ricky -D && mkdir -p /home/ricky/bin && chown -Rf ricky:ricky /home/ricky/bin

USER ricky
WORKDIR /home/ricky/bin

COPY --chown=ricky:ricky --from=builder /go/src/main .

EXPOSE 8080

ENTRYPOINT [ "./main" ]