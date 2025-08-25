# Start from a stable, pre-configured Python base image
FROM mcr.microsoft.com/devcontainers/python:3.10-bullseye

# [Optional] Install additional OS packages.
# This is a good place to add system-level dependencies.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends git vim curl

# Set the working directory in the container
WORKDIR /workspace

# Copy the requirements file and install Python dependencies
# This is done early to leverage Docker's layer caching.
COPY ../requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# This command keeps the container alive.
# The dev container spec requires the container to run continuously.
CMD ["/bin/sh", "-c", "while sleep 1000; do :; done"]
