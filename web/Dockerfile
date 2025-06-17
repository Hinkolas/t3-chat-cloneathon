# Stage 1: install & build
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package manifests & install dependencies
COPY package*.json ./
RUN npm ci

# Copy source & build
COPY . .
RUN npm run build

# Stage 2: run built app
FROM node:22-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

# Install only production deps
COPY package*.json ./
RUN npm ci --only=production

# Copy build output from builder
COPY --from=builder /app/build /app/build

EXPOSE 3000
CMD ["node", "build/index.js"]