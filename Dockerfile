FROM node:16-alpine as node

RUN mkdir /dist
COPY ./internal/api/www /dist
WORKDIR /dist
RUN yarn && yarn build

FROM golang:1.18-alpine AS go

RUN apk update && apk add libreoffice
RUN apk add --no-cache msttcorefonts-installer fontconfig
RUN update-ms-fonts

RUN mkdir /app
COPY . /app
WORKDIR /app
COPY --from=node /dist/public/ internal/api/www/public/

RUN go mod download
RUN go build -o main internal/main.go
CMD ["/app/main"]