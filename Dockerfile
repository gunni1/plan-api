# Stage Build
FROM golang:1.11 as build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./
COPY api ./api
RUN go mod download
RUN go build


# Stage Run
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/plan-api .
RUN chmod +x plan-api
CMD ["./plan-api"]
EXPOSE 8080