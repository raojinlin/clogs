FROM golang:1.19.3 as builder

COPY . /code
WORKDIR /code

RUN go mod tidy && CGO_ENABLED=0 go build -o ./clogs .

FROM node as webui-builder
COPY . /code
WORKDIR /code
RUN npm install && npm run build


FROM alpine

COPY ./template /code/template
COPY --from=builder /code/clogs /code/sbin/clogs
COPY --from=webui-builder /code/build /code

WORKDIR /code
ENV GIN_MODE=release

CMD ["/code/sbin/clogs"]


