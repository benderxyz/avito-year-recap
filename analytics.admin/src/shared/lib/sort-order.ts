export function sortOrderBetween<T extends { sortOrder: number }>(
  ordered: T[],
  movedIndex: number,
): number {
  const prev = ordered[movedIndex - 1];
  const next = ordered[movedIndex + 1];

  if (!prev && next) {
    return next.sortOrder - 1;
  }

  if (prev && !next) {
    return prev.sortOrder + 1;
  }

  if (prev && next) {
    if (next.sortOrder - prev.sortOrder > 1) {
      return Math.floor((prev.sortOrder + next.sortOrder) / 2);
    }

    return prev.sortOrder + 1;
  }

  return 0;
}

export function priorityBetween<T extends { priority: number }>(
  ordered: T[],
  movedIndex: number,
): number {
  const prev = ordered[movedIndex - 1];
  const next = ordered[movedIndex + 1];

  if (!prev && next) {
    return next.priority + 1;
  }

  if (prev && !next) {
    return prev.priority - 1;
  }

  if (prev && next) {
    if (prev.priority - next.priority > 1) {
      return Math.floor((prev.priority + next.priority) / 2);
    }

    return prev.priority - 1;
  }

  return 0;
}
