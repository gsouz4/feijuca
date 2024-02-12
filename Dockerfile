FROM golang:1.22 AS builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /rinha cmd/main.go

FROM scratch

WORKDIR /

COPY --from=builder /rinha /rinha

EXPOSE 5000

CMD [ "/rinha" ]