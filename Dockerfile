FROM rust:1.51 as build
WORKDIR /usr/src

RUN rustup target add x86_64-unknown-linux-musl

RUN USER=root cargo new simple_note_api
WORKDIR /usr/src/simple_note_api
COPY Cargo.toml Cargo.lock ./
RUN cargo build --release

COPY src ./src
RUN cargo install --target x86_64-unknown-linux-musl --path .

FROM scratch

COPY --from=build /usr/local/cargo/bin/simple_note_api .

USER 1000
CMD ["./simple_note_api"]
