import type { SceneDefinition } from '@recap-engine/core';

export type SceneRendererProps<TData = unknown> = {
  scene: SceneDefinition<TData>;
};
