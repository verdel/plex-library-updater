ARG golang_version
ARG distroless_static_version

FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/plex-library-updater .

FROM gcr.io/distroless/static:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39
WORKDIR /
COPY --from=build /app/build/plex-library-updater .
USER 65532:65532
ENTRYPOINT ["/plex-library-updater"]
