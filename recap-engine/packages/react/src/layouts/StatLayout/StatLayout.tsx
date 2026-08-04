import { EMotionPreset, resolveValue } from '@recap-engine/core';
import { Comparison, Eyebrow, SceneActions, Stat, Title } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { StatLayoutProps } from './interface';

export function StatLayout<TData>({ scene }: StatLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();
  const value = resolveValue(scene.value, ctx) ?? 0;
  const unit = resolveValue(scene.unit, ctx);
  const comparison = scene.comparison
    ? {
        template: scene.comparison.template,
        percentile: resolveValue(scene.comparison.percentile, ctx) ?? 0,
      }
    : null;
  const animate = (scene.blockMotion ?? EMotionPreset.CountUp) === EMotionPreset.CountUp;

  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Eyebrow>{resolveValue(scene.eyebrow, ctx)}</Eyebrow>
        <Title>{resolveValue(scene.title, ctx)}</Title>
        <Stat value={value} unit={unit} valueFormat={scene.valueFormat} animate={animate} />
        {comparison ? (
          <Comparison template={comparison.template} percentile={comparison.percentile} />
        ) : null}
      </div>
      <SceneActions actions={scene.actions} />
    </div>
  );
}
