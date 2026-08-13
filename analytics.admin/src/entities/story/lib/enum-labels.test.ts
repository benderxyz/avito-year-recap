import { describe, expect, it } from 'vitest';
import { i18n } from '@/shared/i18n/config';
import { getStorySceneTypeLabel, getStoryVisibilityLabel } from './enum-labels';

describe('story enum labels', () => {
  it('localizes known visibility and sceneType values', () => {
    expect(getStoryVisibilityLabel(i18n.t.bind(i18n), 'private')).toBe('Private');
    expect(getStoryVisibilityLabel(i18n.t.bind(i18n), 'public')).toBe('Public');
    expect(getStorySceneTypeLabel(i18n.t.bind(i18n), 'intro')).toBe('Intro');
    expect(getStorySceneTypeLabel(i18n.t.bind(i18n), 'stat')).toBe('Stat');
  });
});
