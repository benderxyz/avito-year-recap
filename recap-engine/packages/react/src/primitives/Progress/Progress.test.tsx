import { EPlayerPhase } from '@recap-engine/core';
import { renderWithRecap } from '../../test/renderWithRecap';
import { Progress } from './Progress';

describe('Progress', () => {
  it('renders positional fills with a visible minimum for the current scene', () => {
    const { container } = renderWithRecap(<Progress />, {
      player: { index: 1, direction: 1, phase: EPlayerPhase.Active, total: 3 },
      progress: 0,
    });

    const fills = container.querySelectorAll<HTMLElement>('.recap-progress__fill');
    expect(fills).toHaveLength(3);
    expect(fills[0].style.getPropertyValue('--fill')).toBe('1');
    expect(fills[1].style.getPropertyValue('--fill')).toBe('0.08');
    expect(fills[2].style.getPropertyValue('--fill')).toBe('0');
    expect(container.querySelector('.recap-progress')).toHaveAttribute('aria-hidden', 'true');
  });
});
