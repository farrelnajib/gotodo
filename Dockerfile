FROM golang:1.17.5-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /gotodo

EXPOSE 3030

CMD [ "/gotodo" ]