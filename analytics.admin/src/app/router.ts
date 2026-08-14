import { createRouter } from '@tanstack/react-router';
import { routeTree } from '@/app/routeTree.gen';
import { getRouterBasepath } from '@/shared/config/app-base';

const basepath = getRouterBasepath();

export const router = createRouter({
  routeTree,
  ...(basepath ? { basepath } : {}),
});

export type AppRouter = typeof router;
