import {
  EBackgroundSlotName,
  EBackgroundType,
  EImageFit,
  ESceneType,
  type SceneDefinition,
} from '@recap-engine/core';
import { backgroundStyle } from './backgroundStyle';
import { resolveSceneBackground } from './resolveSceneBackground';

const scene = (background?: SceneDefinition['background']): SceneDefinition => ({
  id: 'background',
  type: ESceneType.Custom,
  actions: [],
  background,
});

describe('background resolution', () => {
  it('maps color, gradient and image backgrounds to inline styles', () => {
    expect(backgroundStyle({ type: EBackgroundType.Color, value: '#123456' })).toEqual({
      background: '#123456',
    });
    expect(
      backgroundStyle({
        type: EBackgroundType.Gradient,
        value: 'linear-gradient(red, blue)',
      }),
    ).toEqual({ backgroundImage: 'linear-gradient(red, blue)' });
    expect(
      backgroundStyle({
        type: EBackgroundType.Image,
        src: '/contained.png',
        fit: EImageFit.Contain,
      }),
    ).toEqual({
      backgroundImage: 'url(/contained.png)',
      backgroundSize: 'contain',
      backgroundPosition: 'center',
    });
  });

  it('uses cover by default and returns no style without a background', () => {
    expect(backgroundStyle()).toEqual({});
    expect(backgroundStyle({ type: EBackgroundType.Image, src: '/cover.png' })).toEqual(
      expect.objectContaining({ backgroundSize: 'cover' }),
    );
  });

  it('prefers a scene background over the theme asset', () => {
    expect(
      resolveSceneBackground(
        scene({ type: EBackgroundType.Color, value: 'rebeccapurple' }),
        '/theme.png',
      ),
    ).toEqual({
      useSlotBg: false,
      style: { background: 'rebeccapurple' },
    });
  });

  it('falls back to the theme image', () => {
    expect(resolveSceneBackground(scene(), '/theme.png')).toEqual({
      useSlotBg: false,
      style: {
        backgroundImage: 'url(/theme.png)',
        backgroundPosition: 'center',
        backgroundSize: 'cover',
      },
    });
  });

  it('marks slot backgrounds while retaining the themed base', () => {
    expect(
      resolveSceneBackground(
        scene({
          type: EBackgroundType.Slot,
          name: EBackgroundSlotName.Background,
        }),
        '/theme.png',
      ),
    ).toEqual({
      useSlotBg: true,
      style: expect.objectContaining({ backgroundImage: 'url(/theme.png)' }),
    });
  });
});
