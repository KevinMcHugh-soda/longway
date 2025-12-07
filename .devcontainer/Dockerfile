FROM node:20-bookworm-slim

ARG GO_VERSION=1.22.6

# install OS deps as root
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
         git ca-certificates curl bash; \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz; \
    tar -C /usr/local -xzf /tmp/go.tgz; \
    rm /tmp/go.tgz; \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:${PATH}"

# we are root by default, so this USER root is optional
USER root

# install codex globally as root
RUN npm install -g @openai/codex \
    && chmod +x /usr/local/bin/codex

# ---- still root here ----

ENV XDG_CONFIG_HOME=/home/node/.config
ENV XDG_DATA_HOME=/home/node/.local/share
ENV CODEX_HOME=/home/node/.codex

# create codex + xdg dirs as root for now
RUN mkdir -p "$CODEX_HOME" "$XDG_CONFIG_HOME" "$XDG_DATA_HOME"

# copy auth token to codex home (runs as root; file becomes root:root)
# Note: auth.json should be copied to .devcontainer/ before building
COPY auth.json $CODEX_HOME/auth.json

# 🔑 FIX: hand ownership of everything under /home/node to the node user
RUN chown -R node:node /home/node

# now drop to the built-in non-root user from the node image
USER node
WORKDIR /workspaces/project
