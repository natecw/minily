# Minily URL Shortener

This is minimal URL shortener that serves as a learning project for various technologies.

Current tech stack is Postgres, Redis, Golang and nginx proxy.

## Install

    TODO

## Run the app

    docker compose up --build -d

    docker compose --profile tools run migrate

## Run the tests

    ./run-tests.sh

# REST API

## Get to original URL

### Request

`GET /{short_code}/`

    curl -i http://localhost:8080/abc123/

### Response

    HTTP/1.1 307 Temporary Redirect
    Content-Type: text/html; charset=utf-8
    Location: https://google.com
    Date: Sat, 02 Nov 2024 18:54:17 GMT
    Content-Length: 54

    <a href="https://google.com">Temporary Redirect</a>.


## Create a new short url

### Request

`POST /`

    curl localhost:8080 -i -X POST -H "Content-Type: application/json" -d '{"long_url": "https://yahoo.com","created_by":"me"}'

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Sat, 02 Nov 2024 18:55:51 GMT
    Content-Length: 48

    {"short_url":"c88f320dec138ba5ab0a5f990ff082ba"}
