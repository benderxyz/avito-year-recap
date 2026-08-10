export function getPrivateRecapTitle(year: number | undefined): string {
  return `Итоги ${year ?? ''}`;
}

export function getSharedRecapTitle(displayName: string | undefined): string {
  if (displayName) {
    return `С вами поделилась ${displayName}`;
  }

  return 'С вами поделились итогами';
}
