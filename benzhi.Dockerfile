FROM golang:1.23-bookworm AS build
WORKDIR /src
ENV GOPROXY=off
ENV GOSUMDB=off
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor -o /out/columnstore ./cmd/columnstore

FROM golang:1.23-bookworm
WORKDIR /app
ENV GOPROXY=off
ENV GOSUMDB=off
COPY --from=build /out/columnstore ./columnstore
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
EXPOSE 18080
CMD ["/app/columnstore", "-listen", "0.0.0.0:18080"]
