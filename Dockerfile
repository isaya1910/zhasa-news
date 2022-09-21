FROM golang:1.16-alpine

WORKDIR /app
COPY . .

RUN go build -o main main.go

# run
WORKDI /app
COPY --from=builder /app/main .
EXPOSE 8080

CMD [ "/app/main" ]