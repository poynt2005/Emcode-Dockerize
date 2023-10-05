#!/bin/bash
WASI_SDK_PATH="/opt/wasi-sdk"
WASI_SYSROOT="${WASI_SDK_PATH}/share/wasi-sysroot"
CC="${WASI_SDK_PATH}/bin/clang"  # --sysroot=${WASI_SDK_SYSROOT}
CXX="${WASI_SDK_PATH}/bin/clang++"
LD="${WASI_SDK_PATH}/bin/clang-ld"
NM="${WASI_SDK_PATH}/bin/llvm-nm"
AR="${WASI_SDK_PATH}/bin/llvm-ar"
RANLIB="${WASI_SDK_PATH}/bin/llvm-ranlib"