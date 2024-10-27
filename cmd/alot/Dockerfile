# Build Stage
FROM golang:1.22.5-bullseye AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV GOARCH=amd64 CGO_ENABLED=0 
RUN go build \
    -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o /go/bin/app

# Deploy Stage
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /go/bin/app /app

EXPOSE 8081
USER nonroot:nonroot

CMD ["/app"]
