FROM golang:1.15.6-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/main github.com/KGRC199913/ourdemy_backend/cmd/ourdemy
FROM scratch AS bin
COPY --from=build /bin /
COPY --from=build /src/config /config
EXPOSE 8080:8080
CMD ["./main"]
