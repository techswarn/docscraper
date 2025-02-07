FROM golang:1.23-alpine as builder

ARG SPACES_KEY
ENV SPACES_KEY=${SPACES_KEY}
ARG SPACES_SECRET
ENV SPACES_SECRET=${SPACES_SECRET}
ARG SPACES_ENDPOINT
ENV SPACES_ENDPOINT=${SPACES_ENDPOINT}
ARG FILE
ENV FILE=${FILE}
ARG GO_ENV
ENV GO_ENV=${GO_ENV}

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . ./

RUN go build -v -o bin/docscraper

FROM ubuntu:latest

WORKDIR /workspace

RUN apt-get update && \
    apt-get install -y golang-go git ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/docscraper /usr/local/bin/

CMD ["docscraper"]