import { ESceneType } from '@recap-engine/core';
import { AchievementLayout } from '../AchievementLayout';
import { BlocksLayout } from '../BlocksLayout';
import { CustomLayout } from '../CustomLayout';
import { InsightLayout } from '../InsightLayout';
import { IntroLayout } from '../IntroLayout';
import { OutroLayout } from '../OutroLayout';
import { StatLayout } from '../StatLayout';
import { UpsellLayout } from '../UpsellLayout';
import type { SceneRendererProps } from './interface';

export function SceneRenderer<TData>({ scene }: SceneRendererProps<TData>) {
  switch (scene.type) {
    case ESceneType.Intro:
      return <IntroLayout scene={scene} />;
    case ESceneType.Stat:
      return <StatLayout scene={scene} />;
    case ESceneType.Insight:
      return <InsightLayout scene={scene} />;
    case ESceneType.Upsell:
      return <UpsellLayout scene={scene} />;
    case ESceneType.Achievement:
      return <AchievementLayout scene={scene} />;
    case ESceneType.Outro:
      return <OutroLayout scene={scene} />;
    case ESceneType.Blocks:
      return <BlocksLayout scene={scene} />;
    case ESceneType.Custom:
      return <CustomLayout scene={scene} />;
    default: {
      const _exhaustive: never = scene;
      return <CustomLayout scene={_exhaustive} />;
    }
  }
}
