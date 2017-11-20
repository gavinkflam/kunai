############################################################
# Stage 1: Build libvips and kunai within builder

FROM golang:1.9.2-alpine3.6 as builder

WORKDIR /tmp
ENV \
  LIBVIPS_VERSION_MAJOR=8 \
  LIBVIPS_VERSION_MINOR=5 \
  LIBVIPS_VERSION_PATCH=9

RUN \
  apk add --no-cache --virtual .build-deps \
    gcc g++ make libc-dev \
    curl \
    automake \
    libtool \
    tar \
    gettext && \

  apk add --no-cache --virtual .libdev-deps \
    glib-dev \
    expat-dev \
    libpng-dev \
    libwebp-dev \
    libexif-dev \
    libxml2-dev \
    libjpeg-turbo-dev && \

  LIBVIPS_VERSION=${LIBVIPS_VERSION_MAJOR}.${LIBVIPS_VERSION_MINOR}.${LIBVIPS_VERSION_PATCH} && \

  curl -sL https://github.com/jcupitt/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz | tar -zxv && \
  cd vips-${LIBVIPS_VERSION} && \

  ./configure --without-python --without-gsf && \
  make -j4 && \
  make install && \

  rm -rf /tmp/vips-* && \
  rm -rf /var/cache/apk/* && \
  rm -rf /tmp/vips-*

ENV \
  CPATH=/usr/local/include \
  LIBRARY_PATH=/usr/local/lib \
  PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH

# Copy source codes into container
WORKDIR /go/src/github.com/gavinkflam/kunai
COPY . .

RUN \
  go build \
    -i -a -tags netgo \
    -installsuffix netgo \
    -ldflags='-s -w' \
    -o kunai

############################################################
# Stage 2: assemble production image
FROM alpine:3.6

MAINTAINER Gavin Lam <me@gavin.hk>

# Environment vairables
ENV \
  REFRESHED_AT=2017-11-20 \
  # System variables
  LANG=en_US.UTF-8 \
  HOME=/opt/app/ \
  PATH="/opt/app/bin:${PATH}" \
  # Runtime variables
  LIBRARY_PATH=/usr/local/lib \
  # Application variables
  GIN_MODE=release \
  DIR=/mnt/assets \
  PORT=8080 \
  SECRET_KEY=CHANGEME

RUN \
  apk add --no-cache --virtual .run-deps \
    glib \
    expat \
    libpng \
    libwebp \
    libexif \
    libxml2 \
    libjpeg-turbo

# Work in home directory afterwards
WORKDIR ${HOME}

# Expose port
EXPOSE ${PORT}

# Mount directory to serve
VOLUME ${DIR}

# Copy libvips
COPY --from=builder /usr/local/lib /usr/local/lib

# Copy release files
COPY --from=builder /go/src/github.com/gavinkflam/kunai/kunai /opt/app/bin/kunai

# Set startup command
CMD ["sh", "-c", "kunai"]
