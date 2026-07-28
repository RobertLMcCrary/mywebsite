FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/portfolio ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=build /out/portfolio ./portfolio
COPY --from=build /src/content ./content

ENV LISTEN_ADDR=:8080
ENV POSTS_DIR=/app/content/posts

EXPOSE 8080
ENTRYPOINT ["/app/portfolio"]
