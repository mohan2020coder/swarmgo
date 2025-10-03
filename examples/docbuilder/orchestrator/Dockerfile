FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/orchestrator ./cmd/server

FROM alpine:3.18
COPY --from=build /bin/orchestrator /usr/local/bin/orchestrator
WORKDIR /app
EXPOSE 8080
CMD ["/usr/local/bin/orchestrator"]

