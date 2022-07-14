FROM node:16-alpine as node

RUN mkdir /dist
COPY ./www/package.json /dist/package.json
COPY ./www/yarn.lock /dist/yarn.lock
WORKDIR /dist
RUN yarn
COPY ./www /dist
RUN yarn build

FROM golang:1.18-alpine AS go

# for go tests
RUN CGO_ENABLED=0

# libreoffice for doc to pdf
RUN apk update && apk add libreoffice
RUN apk add --no-cache msttcorefonts-installer fontconfig
RUN update-ms-fonts

RUN mkdir /app
RUN mkdir -p /tmp/commonsyllabi/files
COPY ./tests/samples /app/samples
COPY go.mod /app
COPY go.sum /app
WORKDIR /app
RUN go mod download

COPY cmd /app/cmd
COPY internal /app/internal

COPY --from=node /dist/public/ ./www/public

RUN go build -o bin/api ./cmd/api/main.go
CMD ["/app/bin/api"]