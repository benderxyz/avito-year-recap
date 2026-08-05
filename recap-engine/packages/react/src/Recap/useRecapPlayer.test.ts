import { jest } from '@jest/globals';
import { ERecapEventType } from '@recap-engine/core';
import { act, renderHook } from '@testing-library/react';
import { introScene } from '../test/recapTestUtils';
import { useRecapPlayer } from './useRecapPlayer';

describe('useRecapPlayer', () => {
  const scenes = [introScene('first'), introScene('second'), introScene('third')];

  it('selects the requested initial scene and reports progress', () => {
    const { result } = renderHook(() => useRecapPlayer(scenes, 'second'));

    expect(result.current.player.index).toBe(1);
    expect(result.current.player.total).toBe(3);
    expect(result.current.progress).toBe(0.5);
  });

  it('falls back to the first scene for an unknown initial id', () => {
    const { result } = renderHook(() => useRecapPlayer(scenes, 'missing'));

    expect(result.current.player.index).toBe(0);
  });

  it('navigates within boundaries and emits complete at the end', () => {
    const onEvent = jest.fn();
    const { result } = renderHook(() => useRecapPlayer(scenes, undefined, onEvent));

    act(() => result.current.prev());
    expect(result.current.player.index).toBe(0);

    act(() => result.current.next());
    expect(result.current.player.index).toBe(1);

    act(() => result.current.goTo(2));
    expect(result.current.player.index).toBe(2);
    expect(result.current.player.direction).toBe(1);

    act(() => result.current.next());
    expect(result.current.player.index).toBe(2);
    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });

    act(() => result.current.goTo(0));
    expect(result.current.player.direction).toBe(-1);
  });

  it('updates total and clamps the index when scenes shrink', () => {
    const { result, rerender } = renderHook(
      ({ currentScenes }) => useRecapPlayer(currentScenes, 'third'),
      { initialProps: { currentScenes: scenes } },
    );

    rerender({ currentScenes: [scenes[0]] });

    expect(result.current.player).toMatchObject({ index: 0, total: 1 });
    expect(result.current.progress).toBe(1);
  });

  it('uses complete progress for an empty recap', () => {
    const { result } = renderHook(() => useRecapPlayer([], undefined));

    expect(result.current.player).toMatchObject({ index: 0, total: 0 });
    expect(result.current.progress).toBe(1);
  });
});
