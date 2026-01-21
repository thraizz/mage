# Stage 1: Install dependencies
FROM oven/bun:1.2 AS deps
WORKDIR /app
COPY package.json bun.lockb* bun.lock ./ 
RUN bun install --frozen-lockfile

# Stage 2: Build the application
FROM oven/bun:1.2 AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG VITE_GRPC_SERVER_URL=http://localhost:17171
ENV VITE_GRPC_SERVER_URL=${VITE_GRPC_SERVER_URL}
ENV NODE_ENV=production

RUN bun run build

# Stage 3: Production runtime with nginx
FROM nginx:alpine AS runner

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy built static files from the 'build' directory
# NOTE: Ensure your svelte.config.js uses adapter-static and 'build' is the output dir
COPY --from=builder /app/build /usr/share/nginx/html

# Copy custom nginx configuration
# We use the specific path for the default config override
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:3000/ || exit 1

CMD ["nginx", "-g", "daemon off;"]