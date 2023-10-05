FROM golang:alpine AS builder-base
WORKDIR /progress
WORKDIR /usr/local/go/src/
COPY runner ./runner
RUN cd runner \
    && go build runner.go \
    && chmod +x runner \
    && mv ./runner /progress/runner

FROM ubuntu

ENV CODE_PASSWORD=admin

WORKDIR /progress

RUN apt update
RUN apt install curl git wget python3 make gcc g++ -y 

RUN curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash - \
    && . ~/.bashrc \
    && nvm install node \
    && node -v

RUN curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash \
    && . ~/.wasmedge/env

RUN wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-20/wasi-sdk_20.0_amd64.deb -O wasi-sdk.deb \
    && dpkg -i wasi-sdk.deb; apt install -fy \
    && rm wasi-sdk.deb

RUN curl -fsSL https://code-server.dev/install.sh | bash -

WORKDIR /progress
COPY --from=builder-base /progress/runner ./runner

ENV PATH="${PATH}:/opt/wasi-sdk/bin"

WORKDIR /src

EXPOSE 8080

CMD ["/progress/runner"]