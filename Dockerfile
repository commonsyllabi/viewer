FROM node:16-alpine as node

RUN mkdir /dist
COPY ./internal/www/package.json /dist/package.json
COPY ./internal/www/yarn.lock /dist/yarn.lock
WORKDIR /dist
RUN yarn
COPY ./internal/www /dist
RUN yarn build

FROM golang:1.24-alpine AS go

# for go tests
RUN CGO_ENABLED=0

# pandoc for doc to pdf
RUN apk update && apk add pandoc

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
ENV PORT=80
ENV COSYLL_VIEWER_SAMPLES_DIR=/app/samples
ENV COSYLL_VIEWER_PUBLIC_DIR=/app/internal/www/public
ENV DEBUG=true


RUN go build -o bin/api ./cmd/api/main.go
EXPOSE 80
CMD ["/app/bin/api"]
