FROM golang:1.21.0

WORKDIR /app


COPY . .
RUN go mod download
RUN go build -o ./bin/shortener .

CMD "./bin/shortener"