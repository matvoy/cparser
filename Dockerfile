FROM golang:1.17 as builder

WORKDIR /src
COPY . .
WORKDIR cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /build/cparser .

FROM alpine:latest

WORKDIR /service/cmd
COPY --from=builder /build/cparser .

CMD [ "./cparser" ]