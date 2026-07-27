FROM golang:1.24-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

RUN go install github.com/beego/bee/v2@v2.3.0

COPY . .

EXPOSE 8080

CMD ["bee", "run"]