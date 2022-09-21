FROM golang:1.16-alpine

WORKDIR /app

RUN go mod init zhasa-news
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

EXPOSE 8080
