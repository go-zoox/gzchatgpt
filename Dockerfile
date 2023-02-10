# Builder
FROM whatwewant/builder-go:v1.19-1 as builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux \
  GOARCH=amd64 \
  go build \
  -trimpath \
  -ldflags '-w -s -buildid=' \
  -v -o gzchatgpt

# Server
FROM whatwewant/go:v1.19-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-zoox/gzchatgpt"

ARG VERSION=latest

ENV MODE=production

COPY --from=builder /build/gzchatgpt /bin

ENV VERSION=${VERSION}

CMD gzchatgpt server -c /conf/config.yml
