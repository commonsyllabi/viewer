FROM golang:1.18-alpine AS go

RUN apk update && apk add libreoffice
RUN apk add --no-cache msttcorefonts-installer fontconfig
RUN update-ms-fonts

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o main internal/main.go
CMD ["/app/main"]