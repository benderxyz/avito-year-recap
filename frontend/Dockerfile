FROM node:22-alpine AS build

WORKDIR /app/frontend

RUN corepack enable && corepack prepare pnpm@10.28.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml frontend/.npmrc ./
RUN pnpm install --frozen-lockfile

COPY frontend/ ./

ARG VITE_CARDS_BASE_URL=http://localhost:8081
ARG VITE_RECAP_YEAR=2026
ENV VITE_CARDS_BASE_URL=$VITE_CARDS_BASE_URL
ENV VITE_RECAP_YEAR=$VITE_RECAP_YEAR

RUN pnpm build

FROM node:22-alpine AS runtime

WORKDIR /app/frontend

RUN corepack enable && corepack prepare pnpm@10.28.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml frontend/.npmrc ./
RUN pnpm install --frozen-lockfile --prod=false

COPY --from=build /app/frontend/dist ./dist

EXPOSE 3000

CMD ["pnpm", "preview", "--host", "0.0.0.0", "--port", "3000"]
