FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS builder

ENV GO111MODULE=on CGO_ENABLED=0
WORKDIR /work
ARG TARGETOS TARGETARCH

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/argo-curl-plugin -ldflags="-s -w" .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/argo-curl-plugin /usr/bin/argo-curl-plugin

ARG image_version
ENV RELEASE=$image_version

ENTRYPOINT ["argo-curl-plugin"]
