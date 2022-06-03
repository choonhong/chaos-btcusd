FROM golang:1.17-stretch as build
ENV CGO_ENABLED=1
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the app
RUN go build -a -o /app/chaos-btcusd

# Run the compiled app
CMD ["/app/chaos-btcusd"]
EXPOSE 80
