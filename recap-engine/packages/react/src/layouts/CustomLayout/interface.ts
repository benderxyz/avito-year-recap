import type { CustomScene } from '@recap-engine/core';

export type CustomLayoutProps<TData = unknown> = {
  scene: CustomScene<TData>;
};
