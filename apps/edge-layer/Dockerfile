# Use Rust slim-buster image
FROM rust:slim-buster AS builder

# Set the working directory
WORKDIR /usr/src/app

# Copy the Cargo.toml and Cargo.lock files
COPY Cargo.toml Cargo.lock ./

# Create a dummy main.rs file to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies
RUN cargo build --release && rm -rf src

# Copy the source code
COPY . .

# Build the application
RUN cargo build --release

# Use a minimal image for the final stage
FROM debian:buster-slim

# Set the working directory
WORKDIR /usr/src/app

# Copy the built application from the builder stage
COPY --from=builder /usr/src/app/target/release/edge-layer .

# Add cache mount for builds
VOLUME /usr/src/app/target/release

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./edge-layer"]
