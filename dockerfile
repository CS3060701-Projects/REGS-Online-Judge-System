FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    ninja-build \
    time \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 預設執行指令 (可選)
CMD ["/bin/bash"]