import { ERecapEventType, ESceneActionType, type RecapEvent } from '@recap-engine/react';
import { describe, expect, it } from 'vitest';
import { buildRecapTheme, shouldCloseRecapOnEvent } from '@/entities/recap';

describe('recap helpers', () => {
  it('builds light theme with concrete css vars', () => {
    const theme = buildRecapTheme('light');
    expect(theme.cssVars['--recap-fg']).toBe('#1a1a1a');
    expect(theme.cssVars['--recap-accent']).toBe('#00aaff');
  });

  it('builds dark theme with concrete css vars', () => {
    const theme = buildRecapTheme('dark');
    expect(theme.cssVars['--recap-fg']).toBe('#f5f5f5');
    expect(theme.cssVars['--recap-bg']).toBe('#1f1f1f');
  });

  it('closes on complete event', () => {
    const event: RecapEvent = { type: ERecapEventType.Complete };
    expect(shouldCloseRecapOnEvent(event)).toBe(true);
  });

  it('closes on custom close action', () => {
    const event: RecapEvent = {
      type: ERecapEventType.Action,
      action: { type: ESceneActionType.Custom, id: 'close-recap', label: 'Закрыть' },
      sceneId: 'outro',
    };
    expect(shouldCloseRecapOnEvent(event)).toBe(true);
  });
});
