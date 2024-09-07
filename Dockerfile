FROM golang:1.22

ARG SERVICE

WORKDIR /

ADD go.mod .
COPY . .

RUN go build -o /app cmd/$SERVICE/$SERVICE.go

CMD ["/app"]
