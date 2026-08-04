import {
  type CalloutBlock,
  EMotionPreset,
  ESceneBlockType,
  resolveValue,
  type StatBlock,
  type TextBlock,
} from '@recap-engine/core';
import {
  Callout,
  Comparison,
  Eyebrow,
  SceneActions,
  StaggerText,
  Stat,
  Title,
} from '../../primitives';
import { useSceneCtx } from '../useSceneCtx';
import type { BlocksLayoutProps, BlockViewProps } from './interface';

function BlockView<TData>({ block }: BlockViewProps<TData>) {
  const ctx = useSceneCtx<TData>();

  if (block.type === ESceneBlockType.Stat) {
    const statBlock = block as StatBlock<TData>;
    const value = resolveValue(statBlock.value, ctx) ?? 0;

    return (
      <div className="recap-block">
        <Eyebrow>{resolveValue(statBlock.eyebrow, ctx)}</Eyebrow>
        <Title>{resolveValue(statBlock.title, ctx)}</Title>
        <Stat
          value={value}
          unit={resolveValue(statBlock.unit, ctx)}
          valueFormat={statBlock.valueFormat}
          animate={(statBlock.blockMotion ?? EMotionPreset.CountUp) === EMotionPreset.CountUp}
        />
        {statBlock.comparison ? (
          <Comparison
            template={statBlock.comparison.template}
            percentile={resolveValue(statBlock.comparison.percentile, ctx) ?? 0}
          />
        ) : null}
      </div>
    );
  }

  if (block.type === ESceneBlockType.Text) {
    const textBlock = block as TextBlock<TData>;
    const text = resolveValue(textBlock.text, ctx) ?? '';

    return (
      <StaggerText
        text={text}
        animate={(textBlock.blockMotion ?? EMotionPreset.StaggerText) === EMotionPreset.StaggerText}
        className="recap-block"
      />
    );
  }

  const calloutBlock = block as CalloutBlock<TData>;
  return (
    <div className="recap-block">
      <Callout
        animate={(calloutBlock.blockMotion ?? EMotionPreset.CalloutIn) === EMotionPreset.CalloutIn}
      >
        {resolveValue(calloutBlock.text, ctx)}
      </Callout>
    </div>
  );
}

export function BlocksLayout<TData>({ scene }: BlocksLayoutProps<TData>) {
  return (
    <div className="recap-scene-body">
      <div className="recap-scene-content">
        {scene.blocks.map((block, i) => (
          // Blocks have no stable id in the scene contract.
          // biome-ignore lint/suspicious/noArrayIndexKey: block list has no id field
          <BlockView key={i} block={block} />
        ))}
      </div>
      <SceneActions actions={scene.actions} />
    </div>
  );
}
