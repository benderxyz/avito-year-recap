import { jest } from '@jest/globals';
import {
  createTheme,
  EBackgroundSlotName,
  EBackgroundType,
  EPlayerPhase,
  ERecapEventType,
  ESceneType,
  type SceneDefinition,
} from '@recap-engine/core';
import { fireEvent, screen } from '@testing-library/react';
import { renderWithRecap } from '../../test/renderWithRecap';
import { RecapShell } from './RecapShell';

const contentScene: SceneDefinition = {
  id: 'content',
  type: ESceneType.Intro,
  title: 'Основная сцена',
  actions: [],
};

describe('RecapShell', () => {
  it('composes background, custom slots, theme variables and scene accessibility', () => {
    const slotScene: SceneDefinition = {
      ...contentScene,
      background: {
        type: EBackgroundType.Slot,
        name: EBackgroundSlotName.Background,
      },
    };
    const { container } = renderWithRecap(
      <RecapShell
        className="host-class"
        slots={{
          background: <span>Живой фон</span>,
          header: <header>Особый заголовок</header>,
          footer: <footer>Подвал</footer>,
        }}
      />,
      {
        scenes: [slotScene],
        theme: createTheme({
          colors: { accent: '#ff6600' },
          assets: { background: '/theme.jpg' },
        }),
      },
    );

    expect(screen.getByText('Живой фон')).toBeInTheDocument();
    expect(screen.getByRole('banner')).toHaveTextContent('Особый заголовок');
    expect(screen.getByRole('contentinfo')).toHaveTextContent('Подвал');
    expect(screen.getByRole('region', { name: 'Сцена 1 из 1' })).toHaveAttribute(
      'aria-roledescription',
      'slide',
    );
    const root = container.querySelector<HTMLElement>('.recap-root');
    expect(root).toHaveClass('host-class');
    expect(root?.style.getPropertyValue('--recap-accent')).toBe('#ff6600');
    expect(container.querySelector('.recap-background__base')).toHaveStyle({
      backgroundImage: 'url(/theme.jpg)',
    });
    expect(container.querySelector('.recap-progress')).not.toBeInTheDocument();
  });

  it('renders the default segmented header and reduced-motion presentation', () => {
    renderWithRecap(<RecapShell />, {
      scenes: [contentScene, { id: 'second', type: ESceneType.Outro, title: 'Конец' }],
      reducedMotion: true,
    });

    const slide = screen.getByRole('region', { name: 'Сцена 1 из 2' });
    expect(slide).toHaveClass('recap-motion--fade');
    expect(document.querySelectorAll('.recap-progress__segment')).toHaveLength(2);
  });

  it('supports swipe navigation and can disable gestures', () => {
    const enabled = renderWithRecap(<RecapShell />, { scenes: [contentScene] });
    const root = enabled.container.querySelector('.recap-root') as HTMLElement;

    fireEvent.pointerDown(root, { clientX: 120 });
    fireEvent.pointerUp(root, { clientX: 20 });
    fireEvent.pointerDown(root, { clientX: 20 });
    fireEvent.pointerUp(root, { clientX: 100 });
    expect(enabled.value.next).toHaveBeenCalledTimes(1);
    expect(enabled.value.prev).toHaveBeenCalledTimes(1);

    enabled.unmount();
    const disabled = renderWithRecap(<RecapShell gestures={false} />, {
      scenes: [contentScene],
    });
    const disabledRoot = disabled.container.querySelector('.recap-root') as HTMLElement;
    fireEvent.pointerDown(disabledRoot, { clientX: 100 });
    fireEvent.pointerUp(disabledRoot, { clientX: 0 });
    expect(disabled.value.next).not.toHaveBeenCalled();
  });

  it('supports tap navigation when tapNav is enabled', () => {
    const { value, container } = renderWithRecap(<RecapShell tapNav gestures />, {
      scenes: [contentScene, { id: 'second', type: ESceneType.Outro, title: 'Конец' }],
    });
    const root = container.querySelector('.recap-root') as HTMLElement;
    Object.defineProperty(root, 'getBoundingClientRect', {
      value: () => ({ left: 0, width: 300, top: 0, height: 600 }),
    });

    fireEvent.pointerDown(root, { clientX: 40 });
    fireEvent.pointerUp(root, { clientX: 42 });
    fireEvent.pointerDown(root, { clientX: 260 });
    fireEvent.pointerUp(root, { clientX: 258 });

    expect(value.prev).toHaveBeenCalledTimes(1);
    expect(value.next).toHaveBeenCalledTimes(1);
  });

  it('pauses autoplay while pointer is held when holdToPause is enabled', () => {
    const onAutoplayPausedChange = jest.fn();
    const { container } = renderWithRecap(
      <RecapShell holdToPause autoplay onAutoplayPausedChange={onAutoplayPausedChange} />,
      { scenes: [contentScene] },
    );
    const root = container.querySelector('.recap-root') as HTMLElement;

    fireEvent.pointerDown(root, { clientX: 120, pointerId: 1 });
    fireEvent.pointerUp(root, { clientX: 120, pointerId: 1 });

    expect(onAutoplayPausedChange).toHaveBeenCalledWith(true);
    expect(onAutoplayPausedChange).toHaveBeenCalledWith(false);
  });

  it('handles keyboard navigation and completion at the last scene', () => {
    const onEvent = jest.fn();
    const { value } = renderWithRecap(<RecapShell />, {
      scenes: [contentScene, { id: 'last', type: ESceneType.Outro, title: 'Конец' }],
      player: { index: 1, direction: 1, phase: EPlayerPhase.Active, total: 2 },
      onEvent,
    });

    fireEvent.keyDown(window, { key: 'ArrowRight' });
    fireEvent.keyDown(window, { key: 'ArrowLeft' });
    fireEvent.keyDown(window, { key: ' ' });

    expect(value.next).toHaveBeenCalledTimes(2);
    expect(value.prev).toHaveBeenCalledTimes(1);
    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
  });

  it('renders nothing when the player has no active scene', () => {
    const { container } = renderWithRecap(<RecapShell />, {
      scenes: [],
      player: { index: 0, direction: 1, phase: EPlayerPhase.Active, total: 0 },
    });
    expect(container).toBeEmptyDOMElement();
  });
});
