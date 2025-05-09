FROM golang:1.24-alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go build -o /app/goisolator /app

FROM scratch
COPY --from=backend /app/goisolator .
USER 1000
ENTRYPOINT ["/goisolator"]
