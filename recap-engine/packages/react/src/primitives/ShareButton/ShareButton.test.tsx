import { jest } from '@jest/globals';
import { ERecapEventType, EShareKind, type RecapEvent } from '@recap-engine/core';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRecap } from '../../test/renderWithRecap';
import { ShareButton } from './ShareButton';

describe('ShareButton', () => {
  it('uses the Web Share API and emits share plus action events', async () => {
    const share = jest.fn<() => Promise<void>>().mockResolvedValue();
    Object.defineProperty(navigator, 'share', { configurable: true, value: share });
    const onEvent = jest.fn<(event: RecapEvent) => void>();

    renderWithRecap(
      <ShareButton
        label="Поделиться"
        share={{ kind: EShareKind.Link, title: 'Итоги', text: 'Мой recap', url: '/recap' }}
      />,
      { onEvent },
    );

    await userEvent.click(screen.getByRole('button', { name: 'Поделиться' }));
    await waitFor(() =>
      expect(share).toHaveBeenCalledWith({ title: 'Итоги', text: 'Мой recap', url: '/recap' }),
    );
    expect(onEvent.mock.calls.map(([event]) => event.type)).toEqual([
      ERecapEventType.Share,
      ERecapEventType.Action,
    ]);
  });

  it('copies the URL when native sharing is unavailable', async () => {
    Object.defineProperty(navigator, 'share', { configurable: true, value: undefined });
    const writeText = jest.fn<() => Promise<void>>().mockResolvedValue();
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    renderWithRecap(
      <ShareButton label="Копировать" share={{ kind: EShareKind.Link, url: '/copy-me' }} />,
    );
    await userEvent.click(screen.getByRole('button', { name: 'Копировать' }));

    await waitFor(() => expect(writeText).toHaveBeenCalledWith('/copy-me'));
  });
});
