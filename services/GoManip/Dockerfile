FROM trolllemon/opencv4.11.0 AS gocv-builder 

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl ca-certificates  pkg-config gcc g++ libgtk-3-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev \
    libpng-dev libv4l-dev libxvidcore-dev libx264-dev libtbb-dev libjpeg-dev libtiff-dev libopenexr-dev && \
    rm -rf /var/lib/apt/lists/*

ENV GO_VERSION=1.24.0
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:${PATH}"
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
COPY --from=gocv-builder /usr/local/lib /usr/local/lib
COPY --from=gocv-builder /usr/local/include/opencv4 /usr/local/include/opencv4
COPY --from=gocv-builder /usr/local/share/opencv4 /usr/local/share/opencv4

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
COPY --from=gocv-builder /opt/app/libs /usr/lib/

COPY --from=gocv-builder /app/gomanip /

EXPOSE 8080

ENTRYPOINT ["/gomanip"]
