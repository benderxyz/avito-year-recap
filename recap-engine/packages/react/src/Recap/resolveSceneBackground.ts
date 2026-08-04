import {
  EBackgroundType,
  EImageFit,
  type SceneBackground,
  type SceneDefinition,
} from '@recap-engine/core';
import { backgroundStyle } from './backgroundStyle';

export function resolveSceneBackground<TData>(
  scene: SceneDefinition<TData>,
  themeBackgroundUrl?: string,
) {
  const useSlotBg = scene.background?.type === EBackgroundType.Slot;

  const themeBg: SceneBackground | undefined = themeBackgroundUrl
    ? {
        type: EBackgroundType.Image,
        src: themeBackgroundUrl,
        fit: EImageFit.Cover,
      }
    : undefined;

  const activeBg =
    scene.background && scene.background.type !== EBackgroundType.Slot ? scene.background : themeBg;

  return {
    useSlotBg,
    style: backgroundStyle(activeBg),
  };
}
