FROM golang:1.23.1
WORKDIR /app
EXPOSE 80

COPY . .

RUN go install -mod vendor

ENTRYPOINT ["go", "build", "user-service"]