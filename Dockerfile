FROM golang:1.20.3-alpine
ENV GO111MODULE=on

RUN mkdir /app
WORKDIR /app
ADD . /app
RUN apk add git

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /book-store .
EXPOSE 9011
CMD ["/book-store"]