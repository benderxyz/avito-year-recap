import { EPlayerPhase } from '@recap-engine/core';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRecap } from '../../test/renderWithRecap';
import { Progress } from './Progress';

describe('Progress', () => {
  it('renders completed, current and pending segments with accessible scene progress', () => {
    const { container } = renderWithRecap(<Progress />, {
      player: { index: 1, direction: 1, phase: EPlayerPhase.Active, total: 3 },
    });

    const fills = container.querySelectorAll<HTMLElement>('.recap-progress__fill');
    expect(fills).toHaveLength(3);
    expect(fills[0]).toHaveAttribute('data-state', 'complete');
    expect(fills[1]).toHaveAttribute('data-state', 'current');
    expect(fills[2]).toHaveAttribute('data-state', 'pending');
    expect(screen.getByRole('navigation', { name: 'Навигация по сценам' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Перейти к сцене 2 из 3' })).toHaveAttribute(
      'aria-current',
      'step',
    );
  });

  it('starts the selected scene when a segment is pressed', async () => {
    const { value } = renderWithRecap(<Progress />, {
      autoplay: true,
      player: { index: 0, direction: 1, phase: EPlayerPhase.Active, total: 2 },
    });

    await userEvent.click(screen.getByRole('button', { name: 'Перейти к сцене 2 из 2' }));

    expect(value.goTo).toHaveBeenCalledWith(1);
  });

  it('animates the current segment in sync with autoplay when playback is ready', () => {
    const { container } = renderWithRecap(<Progress />, {
      autoplay: { delayMs: 5000 },
      blockMotionDone: true,
      isAnimating: false,
      player: { index: 0, direction: 1, phase: EPlayerPhase.Active, total: 2 },
      reducedMotion: false,
    });

    const current = container.querySelector<HTMLElement>('[data-state="current"]');
    expect(current).toHaveAttribute('data-animated', 'true');
    expect(current).toHaveAttribute('data-playing', 'true');
    expect(current?.style.getPropertyValue('--recap-progress-duration')).toBe('5000ms');
  });

  it('keeps autoplay progress paused during scene motion', () => {
    const { container } = renderWithRecap(<Progress />, {
      autoplay: true,
      blockMotionDone: true,
      isAnimating: true,
      player: { index: 0, direction: 1, phase: EPlayerPhase.Enter, total: 2 },
      reducedMotion: false,
    });

    expect(container.querySelector('[data-state="current"]')).toHaveAttribute(
      'data-playing',
      'false',
    );
  });

  it('does not render for a single scene', () => {
    const { container } = renderWithRecap(<Progress />, {
      player: { index: 0, direction: 1, phase: EPlayerPhase.Active, total: 1 },
    });

    expect(container.querySelector('.recap-progress')).not.toBeInTheDocument();
  });
});
