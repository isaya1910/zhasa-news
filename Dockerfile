FROM golang:1.16-alpine AS builder

WORKDIR /zhasa-news
COPY . .
ENV GOOGLE_APPLICATION_CREDENTIALS "/app/service-account-file.json"
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /zhasa-news/main .

EXPOSE 8080

CMD [ "/zhasa-news" ]