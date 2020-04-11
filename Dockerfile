FROM ubuntu:16.04

# WIP...
RUN apt update
RUN apt install git libopencv-dev wget libssl-dev -y

RUN wget https://github.com/Kitware/CMake/releases/download/v3.16.2/cmake-3.16.2.tar.gz && \
    tar -xzf cmake-*.tar.gz && cd cmake-* && ./bootstrap && make -j 4 && make install

WORKDIR /build/

RUN git clone --recursive https://github.com/simon987/fastimagehash

WORKDIR /build/fastimagehash

RUN cmake .
RUN make && make install
