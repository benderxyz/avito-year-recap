import { jest } from '@jest/globals';
import {
  createTheme,
  EMotionPreset,
  ERecapEventType,
  ESceneActionType,
  ESceneType,
  type RecapEvent,
  type SceneDefinition,
} from '@recap-engine/core';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Recap } from './Recap';

const theme = createTheme({ colors: { accent: '#ff6600' } });

function flowScenes(): SceneDefinition[] {
  return [
    {
      id: 'intro',
      type: ESceneType.Intro,
      title: 'Ваш 2025',
      subtitle: 'Начинаем',
    },
    {
      id: 'stat',
      type: ESceneType.Stat,
      title: 'Объявления',
      value: 42,
      unit: 'штуки',
      blockMotion: EMotionPreset.None,
      actions: [
        { type: ESceneActionType.Prev, label: 'К началу' },
        { type: ESceneActionType.Next, label: 'К награде' },
      ],
    },
    {
      id: 'badge',
      type: ESceneType.Achievement,
      badgeId: 'champion',
      blockMotion: EMotionPreset.None,
      actions: [{ type: ESceneActionType.GoTo, label: 'Сразу в финал', sceneId: 'outro' }],
    },
    {
      id: 'outro',
      type: ESceneType.Outro,
      title: 'Готово',
      actions: [{ type: ESceneActionType.Custom, label: 'Смотреть снова', id: 'restart' }],
    },
  ];
}

describe('Recap integration flow', () => {
  it('navigates the complete action flow and emits lifecycle/action events', async () => {
    const onEvent = jest.fn<(event: RecapEvent) => void>();
    render(
      <Recap
        theme={theme}
        data={{}}
        scenes={flowScenes()}
        badges={[
          {
            id: 'champion',
            title: 'Чемпион',
            description: 'Лучший результат',
          },
        ]}
        reducedMotion
        onEvent={onEvent}
      />,
    );

    expect(screen.getByRole('region', { name: 'Сцена 1 из 4' })).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Начать' }));
    expect(await screen.findByRole('heading', { name: 'Объявления' })).toBeInTheDocument();
    expect(screen.getByText('42')).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: 'К награде' }));
    expect(await screen.findByRole('heading', { name: 'Чемпион' })).toBeInTheDocument();
    expect(screen.getByText('Лучший результат')).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: 'Сразу в финал' }));
    expect(await screen.findByRole('heading', { name: 'Готово' })).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: 'Смотреть снова' }));
    expect(await screen.findByRole('heading', { name: 'Ваш 2025' })).toBeInTheDocument();

    await waitFor(() => {
      const types = onEvent.mock.calls.map(([event]) => event.type);
      expect(types).toContain(ERecapEventType.SceneEnter);
      expect(types).toContain(ERecapEventType.SceneExit);
      expect(types.filter((type) => type === ERecapEventType.Action)).toHaveLength(4);
    });
    expect(onEvent).toHaveBeenCalledWith(
      expect.objectContaining({
        type: ERecapEventType.SceneExit,
        sceneId: 'outro',
        index: 3,
      }),
    );
  });

  it('supports initial scene selection, keyboard controls and swipe gestures', async () => {
    const onEvent = jest.fn();
    const { container } = render(
      <Recap
        theme={theme}
        data={{}}
        scenes={flowScenes()}
        initialSceneId="stat"
        reducedMotion
        onEvent={onEvent}
      />,
    );

    expect(screen.getByRole('region', { name: 'Сцена 2 из 4' })).toBeInTheDocument();
    fireEvent.keyDown(window, { key: 'ArrowLeft' });
    expect(await screen.findByRole('heading', { name: 'Ваш 2025' })).toBeInTheDocument();

    const root = container.querySelector('.recap-root') as HTMLElement;
    fireEvent.pointerDown(root, { clientX: 100 });
    fireEvent.pointerUp(root, { clientX: 0 });
    expect(await screen.findByRole('heading', { name: 'Объявления' })).toBeInTheDocument();

    fireEvent.keyDown(window, { key: 'ArrowRight' });
    expect(await screen.findByRole('heading', { name: 'Достижение' })).toBeInTheDocument();
    fireEvent.keyDown(window, { key: 'ArrowRight' });
    expect(await screen.findByRole('heading', { name: 'Готово' })).toBeInTheDocument();
    fireEvent.keyDown(window, { key: 'ArrowRight' });
    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
  });

  it('autoplays after motion settles and loops after completion', async () => {
    jest.useFakeTimers();
    const onEvent = jest.fn();
    const scenes: SceneDefinition[] = [
      { id: 'one', type: ESceneType.Intro, title: 'Один', actions: [] },
      { id: 'two', type: ESceneType.Outro, title: 'Два', actions: [] },
    ];

    render(
      <Recap
        theme={theme}
        data={{}}
        scenes={scenes}
        autoplay={{ delayMs: 100 }}
        loop
        reducedMotion
        onEvent={onEvent}
      />,
    );

    await act(async () => {
      jest.advanceTimersByTime(0);
    });
    await act(async () => {
      jest.advanceTimersByTime(100);
    });
    expect(screen.getByRole('heading', { name: 'Два' })).toBeInTheDocument();

    await act(async () => {
      jest.advanceTimersByTime(0);
    });
    await act(async () => {
      jest.advanceTimersByTime(100);
    });
    expect(screen.getByRole('heading', { name: 'Один' })).toBeInTheDocument();
    expect(onEvent).toHaveBeenCalledWith({ type: ERecapEventType.Complete });
    jest.useRealTimers();
  });
});
