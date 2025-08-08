FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o backend

EXPOSE 5000

CMD ["./backend"]