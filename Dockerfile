FROM ubuntu

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install git build-essential libopencv-dev wget libssl-dev -y

RUN wget https://github.com/Kitware/CMake/releases/download/v3.16.2/cmake-3.16.2.tar.gz && \
    tar -xzf cmake-*.tar.gz && cd cmake-* && ./bootstrap && make -j 4 && make install

RUN wget http://fftw.org/fftw-3.3.8.tar.gz && tar -xzf fftw-3.3.8.tar.gz && cd fftw-3.3.8 && ./configure --enable-shared --disable-static --enable-threads --with-combined-threads --enable-portable-binary CFLAGS='-fPIC' && make -j 4 && make install

RUN wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz

WORKDIR /build/

RUN git  clone --recursive  https://github.com/simon987/fastimagehash

WORKDIR /build/fastimagehash

RUN cmake .
RUN make -j 4 && make install

WORKDIR /build/

COPY . /build/imhashdb

WORKDIR /build/imhashdb/cli

RUN PATH=$PATH:/usr/local/go/bin go build .

ENV LD_LIBRARY_PATH /usr/local/lib/

ENTRYPOINT ["/build/imhashdb/cli/cli"]