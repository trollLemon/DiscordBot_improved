FROM debian:bookworm AS opencv-builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential cmake git pkg-config curl ca-certificates \
    libgtk-3-dev libavcodec-dev libavformat-dev libswscale-dev libv4l-dev \
    libxvidcore-dev libx264-dev libjpeg-dev libpng-dev libtiff-dev libopenexr-dev \
    libtbb-dev libeigen3-dev libatlas-base-dev gfortran \
    unzip wget && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt
RUN git clone --branch 4.11.0 --depth 1 https://github.com/opencv/opencv.git && \
    git clone --branch 4.11.0 --depth 1 https://github.com/opencv/opencv_contrib.git

WORKDIR /opt/opencv/build
RUN cmake -D CMAKE_BUILD_TYPE=Release \
          -D CMAKE_INSTALL_PREFIX=/usr/local \
          -D OPENCV_EXTRA_MODULES_PATH=/opt/opencv_contrib/modules \
          -D OPENCV_GENERATE_PKGCONFIG=ON \
          -D BUILD_EXAMPLES=OFF \
          -D BUILD_TESTS=OFF \
          -D BUILD_DOCS=OFF \
          -D BUILD_opencv_python3=OFF \
          .. && \
    make -j$(nproc) && make install && ldconfig


FROM debian:bookworm AS go-builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates  pkg-config gcc g++ libgtk-3-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev \
    libpng-dev libv4l-dev libxvidcore-dev libx264-dev libtbb-dev libjpeg-dev libtiff-dev libopenexr-dev && \
    rm -rf /var/lib/apt/lists/*

ENV GO_VERSION=1.24.0
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:${PATH}"

# Copy OpenCV libraries and includes from the opencv-builder stage
COPY --from=opencv-builder /usr/local/lib /usr/local/lib
COPY --from=opencv-builder /usr/local/include/opencv4 /usr/local/include/opencv4
COPY --from=opencv-builder /usr/local/share/opencv4 /usr/local/share/opencv4

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
ENV CGO_CPPFLAGS="-I/usr/local/include/opencv4"
ENV CGO_LDFLAGS="-L/usr/local/lib"
ENV CGO_ENABLED=1

RUN pkg-config --cflags --libs opencv4

RUN go version
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o gomanip .
RUN ldconfig
RUN mkdir -p /opt/app/libs
RUN ldd /app/gomanip
RUN ldd /app/gomanip | grep "=>" | awk '{print $3}' | xargs -I '{}' sh -c 'if [ -f "{}" ]; then cp "{}" /opt/app/libs/; fi'



FROM golang:tip-bookworm AS tester
COPY --from=opencv-builder /usr/local/lib /usr/local/lib
COPY --from=opencv-builder /usr/local/include/opencv4 /usr/local/include/opencv4
COPY --from=opencv-builder /usr/local/share/opencv4 /usr/local/share/opencv4

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
ENV CGO_CPPFLAGS="-I/usr/local/include/opencv4"
ENV CGO_LDFLAGS="-L/usr/local/lib"
ENV CGO_ENABLED=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates  pkg-config gcc g++ libgtk-3-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev \
    libpng-dev libv4l-dev libxvidcore-dev libx264-dev libtbb-dev libjpeg-dev libtiff-dev libopenexr-dev && \
    rm -rf /var/lib/apt/lists/*

RUN pkg-config --cflags --libs opencv4

RUN go version
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ldconfig

CMD ["go", "test", "-v", "-cover", "-race", "-count=1", "./..."]


FROM gcr.io/distroless/cc-debian12 AS application

ENV LD_LIBRARY_PATH=/usr/local/lib
COPY --from=go-builder /opt/app/libs /usr/lib/

COPY --from=go-builder /app/gomanip /

EXPOSE 8080

ENTRYPOINT ["/gomanip"]
