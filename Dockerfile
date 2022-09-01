FROM alpine:3.14 AS base

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
WORKDIR /app
EXPOSE 6699

FROM golang:1.19-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN go mod tidy
RUN go build -o /bin/repo-stats

FROM base AS runtime
COPY --from=build /bin/repo-stats ./
CMD ["./repo-stats"]
