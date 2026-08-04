import type { RecapContext } from '@recap-engine/core';

export type CustomSceneComponentProps<TData = unknown> = RecapContext<TData> & {
  goNext: () => void;
  goPrev: () => void;
  props?: Record<string, unknown>;
};
