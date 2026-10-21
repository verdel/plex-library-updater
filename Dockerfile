ARG golang_version
ARG distroless_static_version

FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/plex-library-updater .

FROM gcr.io/distroless/static:nonroot@sha256:d62067879de890a1def9da0e7766b7d20810a5e5555cef8506b752437d7a3e88
WORKDIR /
COPY --from=build /app/build/plex-library-updater .
USER 65532:65532
ENTRYPOINT ["/plex-library-updater"]
