FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o /docker-gs-ping
COPY *.go ./

EXPOSE 8080

CMD [ "/zhasa-news" ]