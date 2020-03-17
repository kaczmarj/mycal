FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build .

FROM scratch
COPY --from=builder /app/mycal /mycal
ENTRYPOINT ["/mycal"]
