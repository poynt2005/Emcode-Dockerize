# Emcode-Dockerize

## 動機

由於最近比較常用到 [emscription](https://github.com/emscripten-core/emscripten) 來製作 wasm 的相關程序  
有鑒於 docker 可以快速部署所需環境，並就想結合 [code-server](https://github.com/coder/code-server) 與 emscription 來進行快速開發  
並有了構建這個容器的想法

## 特色

1. 內含實用腳本，有以下功能
   - 可自動安裝 clangd(C/C++ 插件), Prettier(前端插件), MS-Python(python 插件) 等實用插件
   - 可自動設定 format on save 功能
   - 可自動將預設主題設為深色
2. 可直接使用 emscription sdk 工具鏈
3. 結合 vscode

## 使用

直接從 docker hub 拉容器  
再來瀏覽器打開 _你的 server 地址_:_你的 VSCode 埠_ 即可進入 VSCode 介面  
至於運行命令，請參考 runscript.txt
其中：

1. your_code_server_port 就是你想開給 code-server 的 port
2. path/to/your/userdata 用於存放 vscode 用戶資料的卷宗
3. path/to/your/extensions 用於存放 vscode 所安裝之插件的卷宗
4. path/to/your/source_code 由於 emscription 的[官方鏡像](https://hub.docker.com/r/emscripten/emsdk) 是使用 /src 作為工作目錄，因為本容器使用 emscription 的官方鏡像作為基底鏡像，故這邊卷宗的映射也是使用 /src，這邊存放的是你的源碼數據
5. yourpassword 就是 code-server 的密碼，預設 admin

## 備註

1. 原本使用 node.js 腳本做初始設定，現在改為使用 go
2. 原本並無自動設定、多樣卷宗功能，現已加入
