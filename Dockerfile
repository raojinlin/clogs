FROM golang:1.19.3 as builder

COPY . /code
WORKDIR /code

RUN go mod tidy && CGO_ENABLED=0 go build -o ./clogs .


FROM alpine

COPY ./template /code/template
COPY --from=builder /code/clogs /usr/bin/clogs

WORKDIR /code
ENV GIN_MODE=release

CMD ["/usr/bin/clogs"]


