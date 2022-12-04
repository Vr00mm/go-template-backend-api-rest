# Build Image
FROM golang:1.19.3 AS build

# Create application directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY src/go.mod src/go.sum /go/src/app/
RUN go mod download

ADD ./src /go/src/app
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian11:nonroot

LABEL org.opencontainers.artifact.created="from latest ?"
LABEL org.opencontainers.artifact.description="Artefact built in golang 1.19.3"
LABEL org.opencontainers.image.create="current date"
LABEL org.opencontainers.image.authors=["Vr00mm"]
LABEL org.opencontainers.image.url="https://hub.docker.com/vr00mm/application-template"
LABEL org.opencontainers.image.documentation="https://github.com/vr00mm/application-template/wiki"
LABEL org.opencontainers.image.source="https://github.com/vr00mm/application-template/docker"
LABEL org.opencontainers.image.version="0.1.0"
LABEL org.opencontainers.image.revision="ci"
LABEL org.opencontainers.image.vendor="vr00mm"
LABEL org.opencontainers.image.licenses=["NPOSL-3.0"]
LABEL org.opencontainers.image.ref.name="application-template"
LABEL org.opencontainers.image.title="application-template-backend"
LABEL org.opencontainers.image.description=""
LABEL org.opencontainers.image.base.digest="123456"
LABEL org.opencontainers.image.base.name="gcr.io/distroless/base-debian11:nonroot"

COPY --from=build /go/bin/app /

CMD ["/app"]
