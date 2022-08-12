FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . /app

RUN go build -o /zhasa-news

EXPOSE 8080

CMD [ "/zhasa-news" ]
