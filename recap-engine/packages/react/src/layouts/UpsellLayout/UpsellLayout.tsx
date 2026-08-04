import { EMotionPreset, resolveValue } from '@recap-engine/core';
import { Callout, SceneActions, Title } from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { UpsellLayoutProps } from './interface';

export function UpsellLayout<TData>({ scene }: UpsellLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();
  const text = resolveValue(scene.text, ctx);
  const callout = resolveValue(scene.callout, ctx);
  const animate = (scene.blockMotion ?? EMotionPreset.CalloutIn) === EMotionPreset.CalloutIn;

  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Title>{resolveValue(scene.title, ctx)}</Title>
        {text ? <p className="recap-body-text">{text}</p> : null}
        {callout ? <Callout animate={animate}>{callout}</Callout> : null}
      </div>
      <SceneActions actions={scene.actions} />
    </div>
  );
}
