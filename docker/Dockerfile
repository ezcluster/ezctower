# Build the manager binary
FROM  --platform=$BUILDPLATFORM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH


WORKDIR /workspace
# Copy the Go Modules manifests
COPY ../go.mod go.mod
COPY ../go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY ../main.go main.go
COPY ../cmd cmd/
COPY ../internal internal/
COPY ../pkg pkg/

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o tower main.go

FROM ubuntu:focal-20221130
WORKDIR /
COPY --from=builder /workspace/tower /usr/local/bin/tower

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y install --no-install-recommends \
        git \
        python3 \
        python3-pip \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /home/tower && groupadd -g 200 tower &&  useradd --uid 200 -g tower --shell /bin/bash tower && chown -R tower:tower /home/tower

COPY docker/requirements.txt /home/tower/requirements.txt

RUN  pip install --upgrade pip && pip install -r /home/tower/requirements.txt

COPY docker/.profile /home/tower/.profile

USER 200:200

#ENTRYPOINT ["xx"]
CMD ["echo", "Please, provide some commands"]