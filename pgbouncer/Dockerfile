# =========================
# Stage 1: Build
# =========================
FROM alpine:3.19 AS builder

ARG PGBOUNCER_VERSION=1.25.1

RUN apk add --no-cache \
    build-base \
    libevent-dev \
    openssl-dev \
    c-ares-dev \
    postgresql-dev \
    curl \
    ca-certificates \
    pandoc \
    python3

WORKDIR /build

RUN curl -L https://pgbouncer.github.io/downloads/files/${PGBOUNCER_VERSION}/pgbouncer-${PGBOUNCER_VERSION}.tar.gz \
    | tar xz

WORKDIR /build/pgbouncer-${PGBOUNCER_VERSION}

RUN ./configure --prefix=/usr \
    && make -j$(nproc) \
    && make install \
    && strip /usr/bin/pgbouncer

# =========================
# Stage 2: Runtime
# =========================
FROM alpine:3.19

RUN apk add --no-cache \
    libevent \
    openssl \
    c-ares \
    postgresql-libs \
    ca-certificates

RUN addgroup -S pgbouncer && adduser -S pgbouncer -G pgbouncer

COPY --from=builder /usr/bin/pgbouncer /usr/bin/pgbouncer

USER pgbouncer

EXPOSE 6432

ENTRYPOINT ["pgbouncer"]
CMD ["/etc/pgbouncer/pgbouncer.ini"]
