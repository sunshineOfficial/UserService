FROM golang:1.23.1
WORKDIR /app
EXPOSE 80

COPY . .
RUN chmod +x ./start.sh

RUN go install -mod vendor

ENTRYPOINT ["./start.sh"]