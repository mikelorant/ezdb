FROM alpine:3.16

RUN apk add --no-cache \
      mysql-client \
      postgresql-client

COPY ezdb2 ./

ENTRYPOINT ["./ezdb2"]

