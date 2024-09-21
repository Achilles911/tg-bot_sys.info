FROM golang:1.23-alpine

RUN apk update && apk add --no-cache nginx

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o bot .

#CMD ["./bot"] для запуска без nginx

CMD ["sh", "-c", "nginx && ./bot"] 
