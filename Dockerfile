FROM golang:1.22-alpine as dev

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]

FROM scratch as prod

COPY --from=dev /app/main .

ENTRYPOINT ["/main"]