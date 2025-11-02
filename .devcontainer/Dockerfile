# syntax=docker/dockerfile:1
FROM mcr.microsoft.com/devcontainers/go:1-1.23

# Install extra useful tools
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl make netcat-openbsd && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Default workdir set by base image
