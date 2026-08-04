import { createContext, useContext } from 'react';
import type { RecapContextValue, RecapProviderProps } from './interface';

const RecapCtx = createContext<RecapContextValue | null>(null);

export function RecapProvider<TData>({ value, children }: RecapProviderProps<TData>) {
  return <RecapCtx.Provider value={value as RecapContextValue}>{children}</RecapCtx.Provider>;
}

export function useRecap<TData = unknown>(): RecapContextValue<TData> {
  const ctx = useContext(RecapCtx);

  if (!ctx) {
    throw new Error('useRecap must be used within <Recap>');
  }

  return ctx as RecapContextValue<TData>;
}

export type { RecapContextValue } from './interface';
