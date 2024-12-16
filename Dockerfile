FROM rust:1.83
WORKDIR /app
RUN apt update && apt install -y protobuf-compiler
COPY Cargo.toml .
COPY Cargo.lock .
RUN mkdir src \
    && echo "// dummy file" > src/lib.rs \
    && cargo build
COPY src src
COPY build.rs build.rs
COPY service.proto service.proto
RUN cargo build
CMD [ "target/debug/litemq" ]
