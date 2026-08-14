import { describe, expect, it } from 'vitest';
import { priorityBetween, sortOrderBetween } from './sort-order';

describe('sortOrderBetween', () => {
  it('places the item before the first neighbor', () => {
    const ordered = [
      { id: 'moved', sortOrder: 5 },
      { id: 'next', sortOrder: 10 },
    ];

    expect(sortOrderBetween(ordered, 0)).toBe(9);
  });

  it('places the item after the last neighbor', () => {
    const ordered = [
      { id: 'prev', sortOrder: 10 },
      { id: 'moved', sortOrder: 5 },
    ];

    expect(sortOrderBetween(ordered, 1)).toBe(11);
  });

  it('uses the midpoint when neighbors have a gap', () => {
    const ordered = [
      { id: 'prev', sortOrder: 10 },
      { id: 'moved', sortOrder: 1 },
      { id: 'next', sortOrder: 20 },
    ];

    expect(sortOrderBetween(ordered, 1)).toBe(15);
  });

  it('uses prev + 1 when neighbors have no gap', () => {
    const ordered = [
      { id: 'prev', sortOrder: 10 },
      { id: 'moved', sortOrder: 1 },
      { id: 'next', sortOrder: 11 },
    ];

    expect(sortOrderBetween(ordered, 1)).toBe(11);
  });
});

describe('priorityBetween', () => {
  it('places the item above the first neighbor', () => {
    const ordered = [
      { id: 'moved', priority: 5 },
      { id: 'next', priority: 10 },
    ];

    expect(priorityBetween(ordered, 0)).toBe(11);
  });

  it('places the item below the last neighbor', () => {
    const ordered = [
      { id: 'prev', priority: 10 },
      { id: 'moved', priority: 5 },
    ];

    expect(priorityBetween(ordered, 1)).toBe(9);
  });

  it('uses the midpoint when neighbors have a gap', () => {
    const ordered = [
      { id: 'prev', priority: 20 },
      { id: 'moved', priority: 1 },
      { id: 'next', priority: 10 },
    ];

    expect(priorityBetween(ordered, 1)).toBe(15);
  });

  it('uses prev - 1 when neighbors have no gap', () => {
    const ordered = [
      { id: 'prev', priority: 11 },
      { id: 'moved', priority: 1 },
      { id: 'next', priority: 10 },
    ];

    expect(priorityBetween(ordered, 1)).toBe(10);
  });
});
