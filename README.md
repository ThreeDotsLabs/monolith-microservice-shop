# Clean Monolith Shop

Source code for https://threedots.tech/post/microservices-or-monolith-its-detail/ article.

This shop can work both as monolith and microservices. More info you will find in the article.

## Prerequisites

You need **Docker** and **docker-compose** installed.

Everything is running in Docker container, so you don't need golang
either any other lib.

## Running tests

First of all you must run services

    make up


Then you can run all tests by using in another terminal:

    docker-test


If you want to test only monolith version:

    docker-test-monolith

or microservices:

    docker-test-microservices
