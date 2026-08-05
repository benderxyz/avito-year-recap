import {
  EMotionPreset,
  ESceneBlockType,
  ESceneType,
  type SceneDefinition,
} from '@recap-engine/core';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { registerScene } from '../../registry';
import { renderWithRecap } from '../../test/renderWithRecap';
import { SceneRenderer } from './SceneRenderer';

describe('SceneRenderer layouts', () => {
  it('renders intro values from data and its default start action', () => {
    const scene: SceneDefinition<{ name: string }> = {
      id: 'intro',
      type: ESceneType.Intro,
      title: ({ data }) => `Год ${data.name}`,
      subtitle: 'Самое важное',
    };
    renderWithRecap(<SceneRenderer scene={scene} />, { data: { name: 'Максима' } });

    expect(screen.getByRole('heading', { name: 'Год Максима' })).toBeInTheDocument();
    expect(screen.getByText('Самое важное')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Начать' })).toBeInTheDocument();
  });

  it('renders stat, unit and comparison values', () => {
    const scene: SceneDefinition = {
      id: 'stat',
      type: ESceneType.Stat,
      eyebrow: 'Продажи',
      title: 'Ваш результат',
      value: 12500,
      unit: '₽',
      valueFormat: { maximumFractionDigits: 0 },
      comparison: { template: 'Лучше, чем у {{percentile}}% продавцов', percentile: 87 },
      blockMotion: EMotionPreset.None,
      actions: [],
    };
    renderWithRecap(<SceneRenderer scene={scene} />);

    expect(screen.getByText('12 500')).toBeInTheDocument();
    expect(screen.getByText('₽')).toBeInTheDocument();
    expect(screen.getByText('Лучше, чем у 87% продавцов')).toBeInTheDocument();
  });

  it('renders insight defaults and links to another scene', async () => {
    const scene: SceneDefinition = {
      id: 'insight',
      type: ESceneType.Insight,
      title: 'Наблюдение',
      text: 'Первое. Второе!',
      linksTo: 'stat-target',
      blockMotion: EMotionPreset.None,
    };
    const scenes: SceneDefinition[] = [
      scene,
      { id: 'stat-target', type: ESceneType.Stat, value: 1, actions: [] },
    ];
    const { value } = renderWithRecap(<SceneRenderer scene={scene} />, { scenes });

    expect(screen.getByText('Первое.')).toBeInTheDocument();
    expect(screen.getByText('Второе!')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'К цифрам' }));
    expect(value.goTo).toHaveBeenCalledWith(1);
    expect(screen.getByRole('button', { name: 'Дальше' })).toBeInTheDocument();
  });

  it('renders upsell text and optional callout', () => {
    const scene: SceneDefinition = {
      id: 'upsell',
      type: ESceneType.Upsell,
      title: 'Больше возможностей',
      text: 'Попробуйте тариф',
      callout: 'Скидка 20%',
      blockMotion: EMotionPreset.None,
      actions: [],
    };
    renderWithRecap(<SceneRenderer scene={scene} />);

    expect(screen.getByRole('heading', { name: 'Больше возможностей' })).toBeInTheDocument();
    expect(screen.getByText('Попробуйте тариф')).toBeInTheDocument();
    expect(screen.getByText('Скидка 20%')).toBeInTheDocument();
  });

  it('resolves an achievement from the badge catalog and allows scene overrides', () => {
    const scene: SceneDefinition = {
      id: 'achievement',
      type: ESceneType.Achievement,
      badgeId: 'top',
      title: 'Суперпродавец',
      blockMotion: EMotionPreset.None,
      actions: [],
    };
    const { container } = renderWithRecap(<SceneRenderer scene={scene} />, {
      badges: [
        {
          id: 'top',
          title: 'Топ',
          description: 'Вы вошли в десятку',
          icon: '/badge.svg',
        },
      ],
    });

    expect(screen.getByRole('heading', { name: 'Суперпродавец' })).toBeInTheDocument();
    expect(screen.getByText('Вы вошли в десятку')).toBeInTheDocument();
    expect(container.querySelector('.recap-achievement__icon-glyph')).toHaveStyle({
      maskImage: 'url(/badge.svg)',
    });
  });

  it('renders all blocks and supports an empty block collection', () => {
    const scene: SceneDefinition = {
      id: 'blocks',
      type: ESceneType.Blocks,
      blocks: [
        {
          type: ESceneBlockType.Stat,
          eyebrow: 'Всего',
          value: 3,
          unit: 'объявления',
          blockMotion: EMotionPreset.None,
        },
        {
          type: ESceneBlockType.Text,
          text: 'Один вывод.',
          blockMotion: EMotionPreset.None,
        },
        {
          type: ESceneBlockType.Callout,
          text: 'Полезный совет',
          blockMotion: EMotionPreset.None,
        },
      ],
      actions: [],
    };
    const { container, unmount } = renderWithRecap(<SceneRenderer scene={scene} />);

    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('Один вывод.')).toBeInTheDocument();
    expect(screen.getByText('Полезный совет')).toBeInTheDocument();
    expect(container.querySelectorAll('.recap-block')).toHaveLength(3);

    unmount();
    renderWithRecap(
      <SceneRenderer scene={{ id: 'empty', type: ESceneType.Blocks, blocks: [], actions: [] }} />,
    );
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('renders outro defaults', () => {
    const scene: SceneDefinition = {
      id: 'outro',
      type: ESceneType.Outro,
      title: 'До встречи',
      subtitle: 'Это был ваш год',
    };
    renderWithRecap(<SceneRenderer scene={scene} />);

    expect(screen.getByRole('button', { name: 'Поделиться' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'На главную' })).toHaveAttribute('href', '/');
  });
});

describe('custom scenes', () => {
  it('renders an inline scene with context and navigation callbacks', async () => {
    const scene: SceneDefinition<{ label: string }> = {
      id: 'inline',
      type: ESceneType.Custom,
      render: ({ data, goNext }) => (
        <button type="button" onClick={goNext}>
          {data.label}
        </button>
      ),
      actions: [],
    };
    const { value } = renderWithRecap(<SceneRenderer scene={scene} />, {
      data: { label: 'Вперёд из render' },
    });

    await userEvent.click(screen.getByRole('button', { name: 'Вперёд из render' }));
    expect(value.next).toHaveBeenCalledTimes(1);
  });

  it('uses the custom registry and passes scene props', () => {
    registerScene<{ user: string }>('profile-test', ({ data, props, total }) => (
      <p>
        {data.user}: {String(props?.status)} / {total}
      </p>
    ));
    const scene: SceneDefinition<{ user: string }> = {
      id: 'registry',
      type: ESceneType.Custom,
      sceneType: 'profile-test',
      props: { status: 'готов' },
      actions: [],
    };

    renderWithRecap<{ user: string }>(<SceneRenderer scene={scene} />, {
      data: { user: 'Анна' },
      scenes: [scene],
    });
    expect(screen.getByText('Анна: готов / 1')).toBeInTheDocument();
  });

  it('shows a useful fallback for an unknown custom type', () => {
    renderWithRecap(
      <SceneRenderer
        scene={{
          id: 'unknown',
          type: ESceneType.Custom,
          sceneType: 'not-registered',
          actions: [],
        }}
      />,
    );

    expect(
      screen.getByRole('heading', { name: 'Unknown scene: not-registered' }),
    ).toBeInTheDocument();
  });
});
