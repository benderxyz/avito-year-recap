import { describe, expect, it } from 'vitest';
import { PREVIEW_PALETTE_IDS, PREVIEW_PALETTES } from './palettes';

describe('preview palettes', () => {
  it('exposes five named palettes', () => {
    expect(PREVIEW_PALETTE_IDS).toEqual([
      'avitoLight',
      'avitoDark',
      'midnight',
      'sunset',
      'forest',
    ]);
    expect(Object.keys(PREVIEW_PALETTES)).toEqual([...PREVIEW_PALETTE_IDS]);
  });
});
