import { ESceneActionType, resolveValue } from '@recap-engine/core';
import { SceneActions, Subtitle, Title } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { IntroLayoutProps } from './interface';

export function IntroLayout<TData>({ scene }: IntroLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();
  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Title>{resolveValue(scene.title, ctx)}</Title>
        <Subtitle>{resolveValue(scene.subtitle, ctx)}</Subtitle>
      </div>
      <SceneActions actions={scene.actions ?? [{ type: ESceneActionType.Next, label: 'Начать' }]} />
    </div>
  );
}
