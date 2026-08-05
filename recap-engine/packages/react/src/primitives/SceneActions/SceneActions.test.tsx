import { jest } from '@jest/globals';
import { EButtonVariant, ELinkTarget, ERecapEventType, ESceneActionType } from '@recap-engine/core';
import { fireEvent, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderWithRecap } from '../../test/renderWithRecap';
import { SceneActions } from './SceneActions';

describe('SceneActions', () => {
  it('renders default navigation and hides an explicitly empty action list', () => {
    const { unmount, value } = renderWithRecap(<SceneActions />);

    fireEvent.click(screen.getByRole('button', { name: 'Назад' }));
    fireEvent.click(screen.getByRole('button', { name: 'Дальше' }));
    expect(value.prev).toHaveBeenCalledTimes(1);
    expect(value.next).toHaveBeenCalledTimes(1);

    unmount();
    renderWithRecap(<SceneActions actions={[]} />);
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('renders link metadata and emits its action event', async () => {
    const onEvent = jest.fn();
    renderWithRecap(
      <SceneActions
        actions={[
          {
            type: ESceneActionType.Link,
            label: 'Открыть',
            href: 'https://example.test',
            target: ELinkTarget.Blank,
            variant: EButtonVariant.Secondary,
          },
        ]}
      />,
      { onEvent },
    );

    const link = screen.getByRole('link', { name: 'Открыть' });
    expect(link).toHaveAttribute('href', 'https://example.test');
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noreferrer');
    await userEvent.click(link);
    expect(onEvent).toHaveBeenCalledWith(
      expect.objectContaining({ type: ERecapEventType.Action, sceneId: 'scene-1' }),
    );
  });

  it('handles go-to and built-in custom navigation actions', async () => {
    const onEvent = jest.fn();
    const { value } = renderWithRecap(
      <SceneActions
        actions={[
          { type: ESceneActionType.GoTo, label: 'К финалу', sceneId: 'scene-2' },
          { type: ESceneActionType.Custom, label: 'Повторить', id: 'restart' },
          { type: ESceneActionType.Custom, label: 'Вперёд', id: 'next' },
          { type: ESceneActionType.Custom, label: 'Обратно', id: 'prev' },
        ]}
      />,
      { onEvent },
    );

    await userEvent.click(screen.getByRole('button', { name: 'К финалу' }));
    await userEvent.click(screen.getByRole('button', { name: 'Повторить' }));
    await userEvent.click(screen.getByRole('button', { name: 'Вперёд' }));
    await userEvent.click(screen.getByRole('button', { name: 'Обратно' }));

    expect(value.goTo).toHaveBeenNthCalledWith(1, 1);
    expect(value.goTo).toHaveBeenNthCalledWith(2, 0);
    expect(value.next).toHaveBeenCalledTimes(1);
    expect(value.prev).toHaveBeenCalledTimes(1);
    expect(onEvent).toHaveBeenCalledTimes(4);
  });

  it('emits go-to actions even when the target does not exist', async () => {
    const onEvent = jest.fn();
    const { value } = renderWithRecap(
      <SceneActions
        actions={[{ type: ESceneActionType.GoTo, label: 'Нет сцены', sceneId: 'missing' }]}
      />,
      { onEvent },
    );

    await userEvent.click(screen.getByRole('button', { name: 'Нет сцены' }));
    expect(value.goTo).not.toHaveBeenCalled();
    expect(onEvent).toHaveBeenCalledTimes(1);
  });
});
