FROM golang:1.17

RUN mkdir -p /app
WORKDIR /app

# copy the content
COPY . .

# install dependencies
RUN go build -o scootin

# execute
CMD ["./scootin"]