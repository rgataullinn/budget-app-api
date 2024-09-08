# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

ADD . /app
COPY .env /app/.env

RUN go build -o /my-api

EXPOSE 8080

CMD [ "/my-api" ]
