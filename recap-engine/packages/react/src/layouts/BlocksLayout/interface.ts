import type { BlocksScene, SceneBlock } from '@recap-engine/core';

export type BlocksLayoutProps<TData = unknown> = {
  scene: BlocksScene<TData>;
};

export type BlockViewProps<TData = unknown> = {
  block: SceneBlock<TData>;
};
