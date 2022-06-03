FROM golang:1.17-stretch as build
ENV CGO_ENABLED=1
WORKDIR /app

# Install sqlite3
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y sqlite3 libsqlite3-dev
RUN mkdir /db
RUN /usr/bin/sqlite3 /db/test.db

COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the app
RUN go build -a -o /app/chaos-btcusd

# Run the compiled app
CMD ["/app/chaos-btcusd"]
EXPOSE 80
