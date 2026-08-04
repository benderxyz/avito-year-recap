import type { ReactNode } from 'react';
import { SceneActions, Title } from '../../primitives';
import { getRegisteredScene } from '../../registry';
import { useSceneCtx } from '../useSceneCtx';
import type { CustomLayoutProps } from './interface';

export function CustomLayout<TData>({ scene }: CustomLayoutProps<TData>) {
  const ctx = useSceneCtx<TData>();

  if (typeof scene.render === 'function') {
    return (
      <div className="recap-scene-body">
        {
          scene.render({
            ...ctx,
            goNext: ctx.goNext,
            goPrev: ctx.goPrev,
          }) as ReactNode
        }
        {scene.actions !== undefined ? <SceneActions actions={scene.actions} /> : null}
      </div>
    );
  }

  const type = scene.sceneType ?? scene.type;

  const Registered = getRegisteredScene(type);

  if (Registered) {
    return (
      <div className="recap-scene-body">
        <Registered {...ctx} goNext={ctx.goNext} goPrev={ctx.goPrev} props={scene.props} />
        {scene.actions !== undefined ? <SceneActions actions={scene.actions} /> : null}
      </div>
    );
  }

  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        <Title>Unknown scene: {type}</Title>
      </div>
      <SceneActions actions={scene.actions} />
    </div>
  );
}
