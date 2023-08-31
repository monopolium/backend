FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /api ./cmd/api/main.go

EXPOSE 7654

CMD [ "/api" ]