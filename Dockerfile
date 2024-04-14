FROM golang:1.21-alpine

WORKDIR /app

RUN mkdir /config && chmod 777 -R /config

COPY . .

RUN go mod tidy && go build -o app .

ENV TOKEN = ""

CMD sleep 2 && ./app $TOKEN
