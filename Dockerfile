FROM golang:1.18-alpine AS go

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories 
RUN apk add pandoc

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o main internal/main.go
CMD ["/app/main"]