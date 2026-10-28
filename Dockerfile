FROM golang:latest

WORKDIR /kudasai

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . ./

RUN mkdir data

RUN CGO_ENABLED=0 GOOS=linux go build -o /kudasai/kudasai_bot ./cmd/kudasai-bot/

CMD ["/kudasai/kudasai_bot"]


