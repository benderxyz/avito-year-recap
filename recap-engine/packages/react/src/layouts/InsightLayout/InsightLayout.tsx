import { EButtonVariant, EMotionPreset, ESceneActionType, resolveValue } from '@recap-engine/core';
import { SceneActions, StaggerText, Title } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { InsightLayoutProps } from './interface';

export function InsightLayout<TData>({ scene }: InsightLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();
  const text = resolveValue(scene.text, ctx) ?? '';
  const animate = (scene.blockMotion ?? EMotionPreset.StaggerText) === EMotionPreset.StaggerText;

  const actions =
    scene.actions ??
    (scene.linksTo
      ? [
          {
            type: ESceneActionType.GoTo,
            sceneId: scene.linksTo,
            label: 'К цифрам',
            variant: EButtonVariant.Ghost,
          },
          { type: ESceneActionType.Next, label: 'Дальше' },
        ]
      : undefined);

  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Title>{resolveValue(scene.title, ctx)}</Title>
        <StaggerText text={text} animate={animate} className="recap-insight" />
      </div>
      <SceneActions actions={actions} />
    </div>
  );
}
