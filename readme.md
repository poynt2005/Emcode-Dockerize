# Emcode-Dockerize

## 動機
由於最近比較常用到 [emscription](https://github.com/emscripten-core/emscripten) 來製作 wasm 的相關程序  
有鑒於 docker 可以快速部署所需環境，並就想結合 [code-server](https://github.com/coder/code-server) 與 emscription 來進行快速開發  
並有了構建這個容器的想法

## 說明
直接從 docker hub 拉容器  
```bash
docker run --rm -p <your_code_server_port>:8080 -v </path/to/your/volume>:/src -e CODE_PASSWORD="<yourpassword>" -d poynt2005/emcode
```
再來瀏覽器打開 *你的server地址*:*你的VSCode埠* 即可進入 VSCode 介面  
其中： 
1. your_code_server_port 就是你想開給 code-server 的 port  
2. /path/to/your/volume 由於 emscription 的[官方鏡像](https://hub.docker.com/r/emscripten/emsdk) 是使用 /src 作為工作目錄，因為本容器使用 emscription 的官方鏡像作為基底鏡像，故這邊卷宗的映射也是使用 /src，記得進入 code-server 後也要把工作目錄改成 /src  
3. yourpassword 就是 code-server 的密碼，預設 admin  


