FROM ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install wget curl gpg lsb-release -y

RUN wget -O- https://apt.releases.hashicorp.com/gpg | \
    gpg --dearmor | \
    tee /usr/share/keyrings/hashicorp-archive-keyring.gpg

RUN gpg --no-default-keyring \
    --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg \
    --fingerprint

RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
    https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
    tee /etc/apt/sources.list.d/hashicorp.list

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list
RUN apt update

RUN apt install make git -y
RUN apt install terraform -y
RUN apt install python3 python3-pip -y
RUN pip3 install pre-commit
ENV PRE_COMMIT_HOME /workspaces/homelab/.cache

RUN curl -L https://go.dev/dl/go1.20.1.linux-amd64.tar.gz | tar -C /usr/local -xzv
ENV PATH=$PATH:/usr/local/go/bin
RUN go install golang.org/x/tools/gopls@latest 
RUN go install github.com/ramya-rao-a/go-outline@latest
RUN go install github.com/goreleaser/goreleaser@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest