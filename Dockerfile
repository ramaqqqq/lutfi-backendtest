FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go build -o folkatech-try

EXPOSE 7000

CMD ./folkatech-try