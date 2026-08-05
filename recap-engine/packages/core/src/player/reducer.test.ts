import { describe, expect, it } from '@jest/globals';

import { EPlayerActionType, EPlayerPhase, type PlayerState } from '../types/player';
import { createInitialPlayerState, playerReducer } from './reducer';

describe('createInitialPlayerState', () => {
  it.each([
    { total: 5, initialIndex: undefined, expected: 0 },
    { total: 5, initialIndex: 3, expected: 3 },
    { total: 5, initialIndex: -2, expected: 0 },
    { total: 5, initialIndex: 9, expected: 4 },
    { total: 0, initialIndex: 3, expected: 0 },
    { total: -2, initialIndex: 3, expected: 0 },
  ])(
    'clamps index for total=$total and index=$initialIndex',
    ({ total, initialIndex, expected }) => {
      expect(createInitialPlayerState(total, initialIndex)).toEqual({
        index: expected,
        direction: 1,
        phase: EPlayerPhase.Enter,
        total,
      });
    },
  );
});

describe('playerReducer', () => {
  const state = (overrides: Partial<PlayerState> = {}): PlayerState => ({
    index: 1,
    direction: -1,
    phase: EPlayerPhase.Active,
    total: 3,
    ...overrides,
  });

  it('updates total and clamps the current index', () => {
    expect(
      playerReducer(state({ index: 2 }), { type: EPlayerActionType.SetTotal, total: 2 }),
    ).toEqual(state({ index: 1, total: 2 }));
    expect(
      playerReducer(state({ index: 2 }), { type: EPlayerActionType.SetTotal, total: 0 }),
    ).toEqual(state({ index: 0, total: 0 }));
  });

  it('updates phase without changing navigation state', () => {
    expect(
      playerReducer(state(), { type: EPlayerActionType.SetPhase, phase: EPlayerPhase.Exit }),
    ).toEqual(state({ phase: EPlayerPhase.Exit }));
  });

  it('moves next and resets direction and phase', () => {
    expect(playerReducer(state(), { type: EPlayerActionType.Next })).toEqual(
      state({ index: 2, direction: 1, phase: EPlayerPhase.Enter }),
    );
  });

  it('returns the same object when already at the last scene', () => {
    const current = state({ index: 2 });
    expect(playerReducer(current, { type: EPlayerActionType.Next })).toBe(current);
  });

  it('moves back and resets direction and phase', () => {
    expect(playerReducer(state(), { type: EPlayerActionType.Prev })).toEqual(
      state({ index: 0, direction: -1, phase: EPlayerPhase.Enter }),
    );
  });

  it('returns the same object when already at the first scene', () => {
    const current = state({ index: 0 });
    expect(playerReducer(current, { type: EPlayerActionType.Prev })).toBe(current);
  });

  it.each([
    { from: 0, to: 2, direction: 1 as const },
    { from: 2, to: 0, direction: -1 as const },
  ])('jumps from $from to $to in the correct direction', ({ from, to, direction }) => {
    expect(
      playerReducer(state({ index: from }), { type: EPlayerActionType.GoTo, index: to }),
    ).toEqual(state({ index: to, direction, phase: EPlayerPhase.Enter }));
  });

  it.each([-1, 1, 3])('ignores invalid or current go-to index %s', (index) => {
    const current = state();
    expect(playerReducer(current, { type: EPlayerActionType.GoTo, index })).toBe(current);
  });
});
