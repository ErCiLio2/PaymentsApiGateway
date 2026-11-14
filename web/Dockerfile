FROM node:24-alpine AS builder
RUN npm install -g pnpm
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM node:24-alpine
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./
ENV PORT=3000 \
    NODE_ENV=production \
    PROXYPAY_API_URL=http://proxypay:8080
EXPOSE 3000
CMD ["node", "build"]
