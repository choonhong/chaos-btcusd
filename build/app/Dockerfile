FROM golang:1.17-stretch as build
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the app
RUN cd cmd/run && go build -a -o /app/btcusd-api

# Run the compiled app
CMD ["/app/btcusd-api"]