FROM golang:1.23.3

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o scheduler ./cmd/scheduler

EXPOSE 8080
EXPOSE 9090

CMD ["./scheduler"]
