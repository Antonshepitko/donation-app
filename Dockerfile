FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go mod tidy && go mod download
RUN go build -o backend

EXPOSE 5000

CMD ["./backend"]