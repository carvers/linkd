FROM golang:1.27-alpine AS builder

ARG BUILD_ID="unknown"
ARG GIT_SHA="unknown"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X main.BuildID=${BUILD_ID} -X main.GitSHA=${GIT_SHA}" \
    -o linkd .

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/linkd .
EXPOSE 9876
HEALTHCHECK --interval=60s --timeout=1s --start-period=5s CMD curl --fail http://localhost:9876/health || exit 1
CMD ["./linkd"]
