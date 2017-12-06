############################################################
# Stage 1: Build libvips and kunai within builder

FROM golang:1.9.2-alpine3.6 as builder

# Derived from PyYoshi/alpine-libvips
# https://github.com/PyYoshi/alpine-libvips/blob/master/Dockerfile

WORKDIR /tmp

# libvips version
ENV \
  LIBVIPS_VERSION_MAJOR=8 \
  LIBVIPS_VERSION_MINOR=5 \
  LIBVIPS_VERSION_PATCH=9

RUN \
  apk add --no-cache --virtual .build-deps \
    # Used for download and extract libvips source tarball
    curl tar \
    # Dependencies required to build libvips
    gcc g++ make libc-dev \
    automake \
    libtool \
    gettext && \

  # Development files to build vips
  apk add --no-cache --virtual .libdev-deps \
    glib-dev \
    expat-dev \
    libpng-dev \
    libwebp-dev \
    libexif-dev \
    libxml2-dev \
    libjpeg-turbo-dev && \

  # libvips version string
  LIBVIPS_VERSION=${LIBVIPS_VERSION_MAJOR}.${LIBVIPS_VERSION_MINOR}.${LIBVIPS_VERSION_PATCH} && \

  # Download and extract source tarball
  curl -sL https://github.com/jcupitt/libvips/releases/download/v${LIBVIPS_VERSION}/vips-${LIBVIPS_VERSION}.tar.gz | tar -zxv && \
  cd vips-${LIBVIPS_VERSION} && \

  # Build libvips from source
  ./configure --without-python --without-gsf && \
  make -j4 && \
  make install && \

  # Remove caches and temporary files
  rm -rf /tmp/vips-* && \
  rm -rf /var/cache/apk/* && \
  rm -rf /tmp/vips-*

# Specify locations of development files
ENV \
  CPATH=/usr/local/include \
  LIBRARY_PATH=/usr/local/lib \
  PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH

# Copy source codes into container
WORKDIR /go/src/github.com/gavinkflam/kunai
COPY . .

# Build kunai
RUN \
  go build \
    # Install dependencies
    -i \
    # Force rebuild
    -a \
    # Go tool link arguments
    -ldflags=' \
      # Strip debug information to reduce size
      -s -w' \
    # Output
    -o kunai

############################################################
# Stage 2: Assemble production image
FROM alpine:3.6

MAINTAINER Gavin Lam <me@gavin.hk>

# Environment vairables
ENV \
  REFRESHED_AT=2017-12-06 \
  # System variables
  LANG=en_US.UTF-8 \
  HOME=/opt/app/ \
  PATH="/opt/app/bin:${PATH}" \
  # libvips development files
  LIBRARY_PATH=/usr/local/lib \
  # Application variables
  GIN_MODE=release \
  DIR=/mnt/assets \
  HOST=http://localhost \
  PORT=8080 \
  SECRET_KEY=CHANGEME \
  TOKEN=

# Runtime dependencies for libvips
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

# Mount assets to serve
VOLUME ${DIR}

# Copy libvips
COPY --from=builder /usr/local/lib /usr/local/lib

# Copy kunai binary
COPY --from=builder /go/src/github.com/gavinkflam/kunai/kunai /opt/app/bin/kunai

# Set startup command
CMD ["sh", "-c", "kunai"]
