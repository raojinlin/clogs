FROM golang:1.19.3

COPY . /code
WORKDIR /code
RUN go mod tidy && go build -o cri-logs-term-serve .
ENV GIN_MODE=release
CMD ["./cri-logs-term-serve"]


