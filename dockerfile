FROM golang:alpine AS builder-base
WORKDIR /progress
WORKDIR /usr/local/go/src/
COPY runner ./runner
RUN cd runner \
    && go build runner.go \
    && chmod +x runner \
    && mv ./runner /progress/runner

FROM emscripten/emsdk

ENV CODE_PASSWORD=admin

WORKDIR /progress

RUN apt update
RUN apt install curl -y 

RUN curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash - \
    && . ~/.bashrc \
    && nvm install node \
    && node -v

RUN curl -fsSL https://code-server.dev/install.sh | bash -

# 加入現在很潮的 wasm 運行時
RUN curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash \
    && . ~/.wasmedge/env

WORKDIR /progress
COPY --from=builder-base /progress/runner ./runner

WORKDIR /src

EXPOSE 8080

CMD ["/progress/runner"]