import type { ComponentType } from 'react';
import type { CustomSceneComponentProps } from './interface';

const registry = new Map<string, ComponentType<CustomSceneComponentProps>>();

export function registerScene<TData = unknown>(
  type: string,
  component: ComponentType<CustomSceneComponentProps<TData>>,
): void {
  registry.set(type, component as ComponentType<CustomSceneComponentProps>);
}

export function getRegisteredScene(
  type: string,
): ComponentType<CustomSceneComponentProps> | undefined {
  return registry.get(type);
}

export type { CustomSceneComponentProps } from './interface';
