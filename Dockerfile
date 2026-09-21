FROM golang:1.26.6-trixie AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p content && go build -o enchantech ./cmd/api

FROM debian:trixie-slim

# The application reads the weather over HTTPS, so the runtime image needs root
# certificates. The slim image ships without them, and without this the forecast
# call fails quietly and every visitor sees the default sky.
RUN apt-get update \
    && apt-get install --yes --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/enchantech .
COPY --from=build /app/content ./content

EXPOSE 8080
CMD ["./enchantech"]
