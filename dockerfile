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

WORKDIR /progress
COPY --from=builder-base /progress/runner ./runner

WORKDIR /src

EXPOSE 8080

CMD ["/progress/runner"]