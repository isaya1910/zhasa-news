FROM golang:1.16-alpine

WORKDIR /app

COPY go.sum ./
RUN go mod init zhasa-news
COPY go.mod ./
RUN go mod download

COPY *.go ./

EXPOSE 8080

CMD [ "/zhasa-news" ]