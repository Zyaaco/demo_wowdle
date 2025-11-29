# Multi-stage Dockerfile for building a Vite + React app and serving it with nginx
# Build stage: use a lightweight Node image to install deps and build the app
FROM node:24-alpine AS builder

WORKDIR /app

# Install dependencies (copy package files first to leverage Docker cache)
COPY package.json package-lock.json ./
RUN npm ci --silent

# Copy the rest of the app and run the production build
COPY . .
RUN npm run build

# Runtime stage: nginx serving static files
FROM nginx:stable-alpine AS runner

# Copy built assets from the builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Replace default nginx config with one suitable for single-page apps (SPA)
# - fallback to index.html for client-side routing
# - long caching for static assets
RUN rm -f /etc/nginx/conf.d/default.conf
RUN cat > /etc/nginx/conf.d/default.conf <<'EOF'
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # SPA fallback: if file not found, serve index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets aggressively
    location ~* \.(?:css|js|mjs|woff2?|ttf|eot|otf|svg|png|jpe?g|gif|ico)$ {
        expires 30d;
        add_header Cache-Control "public, max-age=2592000, immutable";
    }

    # Optional: serve robots.txt and favicon with short caching
    location = /robots.txt { expires 1h; }
    location = /favicon.ico { expires 7d; }

    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;
}
EOF

EXPOSE 80

# Run nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
