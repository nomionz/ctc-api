FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go mod verify
RUN go build -o ctc-api .

FROM golang:1.20-alpine 

COPY --from=builder /app/ctc-api /app/ctc-api

EXPOSE 8080
CMD ["/app/ctc-api"]