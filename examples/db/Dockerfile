FROM golang:1.15.4-alpine3.12 AS builder

RUN apk update && apk add git

WORKDIR "/"
RUN ["git", "clone", "https://github.com/djordjev/pg-mig"]

WORKDIR "/pg-mig"
RUN ["go",  "build",  "-o", "./build/pg-mig", "./cmd/pg-mig/main.go"]

FROM postgres:13.0-alpine

ENV PG_MIG_CONFIG_DIR=/usr/pg-mig/
ENV PATH="/usr/pg-mig:${PATH}"

COPY ./workspace/ /usr/pg-mig/workspace/

COPY --from=builder /pg-mig/build/pg-mig /usr/pg-mig/pg-mig

RUN chmod -R 777 /usr/pg-mig

EXPOSE 5432