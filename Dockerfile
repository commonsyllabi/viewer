FROM node:16-alpine as node

RUN mkdir /dist
COPY ./www /dist
WORKDIR /dist
RUN yarn && yarn build

FROM golang:1.18-alpine AS go
# for go tests
RUN CGO_ENABLED=0
RUN apk update && apk add libreoffice
RUN apk add --no-cache msttcorefonts-installer fontconfig
RUN update-ms-fonts

RUN mkdir /app
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY go.mod /app
COPY go.sum /app

COPY internal /app/internal

WORKDIR /app
COPY --from=node /dist/public/ ./www/public
RUN go mod download
RUN go build -o bin/api ./cmd/api/main.go
CMD ["/app/bin/api"]