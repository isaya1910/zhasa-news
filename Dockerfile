FROM golang:1.16-alpine

WORKDIR /app
COPY . .

RUN go build -o main main.go

COPY *.go ./

EXPOSE 8080

CMD [ "/zhasa-news" ]