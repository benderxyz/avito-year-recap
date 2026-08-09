FROM node:22-alpine AS build

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@10 --activate

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/core/package.json packages/core/
COPY packages/react/package.json packages/react/
COPY docs/package.json docs/

RUN pnpm install --frozen-lockfile

COPY packages/core packages/core
COPY packages/react packages/react
COPY docs docs

ARG DOCS_URL=https://recaps.hakolr.dev:3001
ENV DOCS_URL=$DOCS_URL

RUN pnpm docs:build

FROM nginx:alpine AS runtime

COPY --from=build /app/docs/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
