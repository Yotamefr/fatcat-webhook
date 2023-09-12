FROM golang:1.21.1-alpine3.18


WORKDIR /app
COPY . /app

ENV GIN_MODE=release
RUN go build
RUN find . ! -name m ! -name LICENSE ! -name . ! -name .. -exec rm -rf {} +

ENTRYPOINT [ "./m" ]
