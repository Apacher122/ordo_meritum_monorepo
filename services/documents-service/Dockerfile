# syntax=docker/dockerfile:1.4

# --- Base layer from prebuilt image with LaTeX ---
FROM ordo_meritum-app:latest AS latex-base
WORKDIR /usr/src/app

# --- Dependencies layer: reuse node_modules if present ---
FROM latex-base AS deps
WORKDIR /usr/src/app

# Copy package.json & package-lock.json first for caching
COPY package*.json ./

# Copy node_modules from base image if present
# COPY ./node_modules ./node_modules

# Install missing packages verbosely and offline if possible
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=cache,target=/usr/src/app/node_modules \
    npm install --prefer-offline --verbose

RUN npm install kafkajs

# --- Build layer ---
FROM deps AS builder
WORKDIR /usr/src/app
COPY . .

# --- Runtime layer ---
FROM latex-base AS runtime
WORKDIR /usr/src/app

# Copy package files
COPY package*.json ./

# Copy node_modules from deps (already installed)
# COPY --from=deps /usr/src/app/node_modules ./node_modules

# Copy built app
COPY --from=builder /usr/src/app ./

# Start with the wrapper script
CMD ["npx", "tsx", "src/index.ts"]
