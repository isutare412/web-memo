# Build the app binary
FROM golang:1.23 as builder

WORKDIR /app

# Copy and install go module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o server ./cmd/...

FROM --platform=linux/amd64 gcr.io/distroless/static:debug

WORKDIR /app
COPY --from=builder /app/server .
COPY config.yaml configs/

ENTRYPOINT [ "./server" ]
CMD [ "-configs", "configs" ]
