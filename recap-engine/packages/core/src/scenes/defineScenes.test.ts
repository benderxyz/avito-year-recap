import { describe, expect, it } from '@jest/globals';

import { ESceneType, type SceneDefinition } from '../types/scenes';
import { defineScenes } from './defineScenes';

describe('defineScenes', () => {
  it('returns the same scene array when ids are unique', () => {
    const scenes: SceneDefinition[] = [
      { id: 'first', type: ESceneType.Custom },
      { id: 'second', type: ESceneType.Custom, sceneType: 'chart' },
    ];

    expect(defineScenes(scenes)).toBe(scenes);
  });

  it('allows values that are similar but not identical', () => {
    expect(() =>
      defineScenes([
        { id: 'scene', type: ESceneType.Custom },
        { id: 'Scene', type: ESceneType.Custom },
        { id: ' scene ', type: ESceneType.Custom },
      ]),
    ).not.toThrow();
  });

  it('rejects duplicate ids and identifies the duplicate', () => {
    expect(() =>
      defineScenes([
        { id: 'duplicate', type: ESceneType.Custom },
        { id: 'other', type: ESceneType.Custom },
        { id: 'duplicate', type: ESceneType.Custom },
      ]),
    ).toThrow('Duplicate scene id: duplicate');
  });
});
