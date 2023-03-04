# TODO connect to database
FROM golang:1.19
MAINTAINER my-fake-mail@gmail.com
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]

