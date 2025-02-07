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
#RUN IMAGE
FROM ubuntu:latest

WORKDIR /workspace

RUN apt-get update && \
    apt-get install -y golang-go git ca-certificates figlet && \
    rm -rf /var/lib/apt/lists/*

# Add crontab file in the cron directory
ADD crontab /etc/cron.d/hello-cron

# Give execution rights on the cron job
RUN chmod 0644 /etc/cron.d/hello-cron

# Create the log file to be able to run tail
RUN touch /var/log/cron.log

#Install Cron
RUN apt-get update
RUN apt-get -y install cron

COPY --from=builder /app/bin/docscraper /usr/local/bin/

# Run the command on container startup
CMD cron && tail -f /var/log/cron.log

