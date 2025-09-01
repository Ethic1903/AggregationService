FROM golang:1.24-alpine AS build

WORKDIR /AggregationService

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/app

FROM alpine:latest
WORKDIR /root/

COPY --from=build /AggregationService/app /AggregationService/app
COPY .env .env

EXPOSE 7071

CMD ["/AggregationService/app"]