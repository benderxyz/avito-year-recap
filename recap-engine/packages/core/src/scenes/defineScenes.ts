import type { SceneDefinition } from '../types/scenes';

export function defineScenes<D>(scenes: SceneDefinition<D>[]): SceneDefinition<D>[] {
  const ids = new Set<string>();
  for (const scene of scenes) {
    if (ids.has(scene.id)) {
      throw new Error(`Duplicate scene id: ${scene.id}`);
    }
    ids.add(scene.id);
  }
  return scenes;
}
