FROM golang:1.17

WORKDIR /src
RUN apt-get update \
		&& apt-get install -y liblzma-dev \
		&& apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY . ./
RUN make

EXPOSE 1633 1634 1635
