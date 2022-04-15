FROM rustlang/rust:nightly AS builder
WORKDIR /opt
COPY . /opt
RUN cargo build --release

FROM debian:stable-slim AS final
WORKDIR /opt
COPY --from=builder /opt/target/release/nitrixme /opt/nitrixme
COPY --from=builder /opt/templates /opt/templates
COPY --from=builder /opt/static /opt/static
EXPOSE 8000
CMD ["/opt/nitrixme"]