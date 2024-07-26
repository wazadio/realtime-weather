FROM golang:latest AS builder

WORKDIR $GOPATH/src/i3
COPY . ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /main main.go

FROM alpine:latest
ENV TZ=Asia/Jakarta
WORKDIR /app
COPY --from=builder /main ./
COPY .env ./
EXPOSE 9000
ENTRYPOINT [ "/app/main" ]