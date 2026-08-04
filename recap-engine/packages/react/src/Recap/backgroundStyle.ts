import { EBackgroundType, EImageFit, type SceneBackground } from '@recap-engine/core';
import type { CSSProperties } from 'react';

export function backgroundStyle(bg?: SceneBackground): CSSProperties {
  if (!bg) return {};
  switch (bg.type) {
    case EBackgroundType.Color:
      return { background: bg.value };
    case EBackgroundType.Gradient:
      return { backgroundImage: bg.value };
    case EBackgroundType.Image:
      return {
        backgroundImage: `url(${bg.src})`,
        backgroundSize: bg.fit ?? EImageFit.Cover,
        backgroundPosition: 'center',
      };
    default:
      return {};
  }
}
