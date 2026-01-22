FROM --platform=linux/amd64 node:24-slim AS base

# Install dependencies only when needed
FROM base AS deps

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Install dependencies based on the preferred package manager
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Rebuild the source code only when needed
FROM base AS builder

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY src src
COPY static static
COPY package.json \
    pnpm-lock.yaml \
    postcss.config.js \
    tailwind.config.js \
    svelte.config.js \
    tsconfig.json \
    vite.config.ts \
    .npmrc \
    ./
RUN pnpm run build

# Production image, copy all the files and run
FROM base AS runner

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app
COPY package.json pnpm-lock.yaml ./
COPY --from=builder /app/build ./server
RUN pnpm install --prod --frozen-lockfile

ENV HOST=0.0.0.0
ENV PORT=80

CMD ["node", "server"]
