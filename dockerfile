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

COPY startapp.js ./

WORKDIR /src

EXPOSE 8080

CMD ["node", "/progress/startapp.js"]