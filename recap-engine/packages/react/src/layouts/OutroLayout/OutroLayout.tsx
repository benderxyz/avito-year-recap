import { EButtonVariant, ESceneActionType, EShareKind, resolveValue } from '@recap-engine/core';
import { SceneActions, Subtitle, Title } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { OutroLayoutProps } from './interface';

export function OutroLayout<TData>({ scene }: OutroLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();
  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Title>{resolveValue(scene.title, ctx)}</Title>
        <Subtitle>{resolveValue(scene.subtitle, ctx)}</Subtitle>
      </div>
      <SceneActions
        actions={
          scene.actions ?? [
            {
              type: ESceneActionType.Share,
              label: 'Поделиться',
              share: { kind: EShareKind.Link },
            },
            {
              type: ESceneActionType.Link,
              label: 'На главную',
              href: '/',
              variant: EButtonVariant.Primary,
            },
          ]
        }
      />
    </div>
  );
}
