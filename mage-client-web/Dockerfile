# Multi-stage Docker build for MAGE Svelte Web Client

# Stage 1: Install dependencies
FROM oven/bun:1.2 AS deps

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Stage 2: Build the application
FROM oven/bun:1.2 AS builder

WORKDIR /app

# Copy node_modules from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy source code
COPY . .

# Build arguments for API URL (can be overridden at build time)
ARG VITE_GRPC_SERVER_URL=http://localhost:17171
ENV VITE_GRPC_SERVER_URL=${VITE_GRPC_SERVER_URL}

# Set production mode for adapter-node
ENV NODE_ENV=production

# Build the application
RUN bun run build

# Stage 3: Production runtime with Node.js
FROM node:22-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 svelte

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

# Change ownership
RUN chown -R svelte:nodejs /app

# Switch to non-root user
USER svelte

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:3000 || exit 1

# Run the application using Node.js
CMD ["node", "build"]

