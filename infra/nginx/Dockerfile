FROM node:22-alpine AS admin-build

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@10.28.0 --activate

COPY analytics.admin/package.json analytics.admin/pnpm-lock.yaml analytics.admin/.npmrc ./
RUN pnpm install --frozen-lockfile

COPY analytics.admin/ ./

ARG VITE_BASE_PATH=/admin/
ARG VITE_API_BASE_URL=
ARG VITE_ADMIN_API_TOKEN
ENV VITE_BASE_PATH=$VITE_BASE_PATH
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_ADMIN_API_TOKEN=$VITE_ADMIN_API_TOKEN

RUN pnpm build

FROM nginx:1.27-alpine

COPY --from=admin-build /app/dist /var/www/admin

EXPOSE 80 443
