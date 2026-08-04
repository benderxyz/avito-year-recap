import { EMotionPreset, resolveValue } from '@recap-engine/core';
import { AchievementBadge, SceneActions } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { AchievementLayoutProps } from './interface';

export function AchievementLayout<TData>({ scene }: AchievementLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();

  const badge = scene.badgeId ? ctx.badges.find((b) => b.id === scene.badgeId) : undefined;

  const title = resolveValue(scene.title, ctx) ?? badge?.title ?? 'Достижение';
  const description = resolveValue(scene.description, ctx) ?? badge?.description;
  const icon = scene.icon ?? badge?.icon;
  const animate = (scene.blockMotion ?? EMotionPreset.BadgePop) === EMotionPreset.BadgePop;

  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content recap-scene-content--center">
        <AchievementBadge title={title} description={description} icon={icon} animate={animate} />
      </div>
      <SceneActions actions={scene.actions} />
    </div>
  );
}
