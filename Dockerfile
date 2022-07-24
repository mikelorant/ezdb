FROM alpine:3.16 as base
RUN ln -s /var/cache/apk /etc/apk/cache
RUN --mount=type=cache,id=apk,target=/var/cache/apk \
    apk add --no-cache \
      mysql-client \
      postgresql-client
WORKDIR /root

FROM golang:1.18-alpine as dependencies
WORKDIR /usr/src/app
COPY go.* ./
RUN --mount=target=. \
    --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download -x

FROM dependencies as build
COPY . ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go build -o ezdb

FROM base as release
RUN apk add --no-cache \
      mysql-client \
      postgresql-client
COPY --from=build /usr/src/app/ezdb /usr/local/bin/

ENTRYPOINT ["ezdb"]
