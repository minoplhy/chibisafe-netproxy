FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 4040

CMD ["go", "run", "main.go"]
