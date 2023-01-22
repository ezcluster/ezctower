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
        curl \
        wget \
        iputils-ping \
        ncat \
        iotop \
        iftop \
        htop \
        mlocate \
        net-tools \
        sudo \
        unzip \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /home/tower && groupadd -g 200 tower &&  useradd --uid 200 -g tower --shell /bin/bash tower && chown -R tower:tower /home/tower

RUN cd /tmp && wget https://releases.hashicorp.com/terraform/1.3.7/terraform_1.3.7_linux_amd64.zip && \
    unzip terraform_1.3.7_linux_amd64.zip && mv /tmp/terraform /usr/local/bin/terraform

RUN curl -L https://github.com/tsl0922/ttyd/releases/download/1.7.3/ttyd.i686 -o /usr/local/bin/ttyd && chmod +x /usr/local/bin/ttyd

COPY docker/requirements.txt /home/tower/requirements.txt

RUN  pip install --upgrade pip && pip install -r /home/tower/requirements.txt

COPY docker/tower.sudoer /etc/sudoers.d/tower
COPY docker/.profile /home/tower/.profile


RUN mkdir -p /home/tower/ezc && cd /home/tower/ezc && \
    git clone --branch substenv --single-branch https://github.com/ezcluster/ezcluster.git && \
    git clone https://github.com/ezcluster/ezcbase.git && \
    git clone https://github.com/ezcluster/ezcmiddle.git && \
    git clone https://github.com/ezcluster/ezck8s.git

RUN mkdir -p /home/tower/nih && cd /home/tower/nih && \
    git clone -b 'v2.18.1' --single-branch --depth 1 https://github.com/kubernetes-sigs/kubespray.git

USER 200:200

#ENTRYPOINT ["xx"]
CMD ["ttyd", "-d", "4", "-p", "7681", "/bin/bash", "-l"]

