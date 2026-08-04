import type { ResolvedTheme } from './theme';

export type PluralForms = {
  one: string;
  few: string;
  many: string;
};

export type Formatters = {
  number: (n: number, opts?: Intl.NumberFormatOptions) => string;
  currency: (n: number, currency?: string) => string;
  plural: (n: number, forms: PluralForms) => string;
};

export type RecapContext<TData = unknown> = {
  data: TData;
  theme: ResolvedTheme;
  index: number;
  total: number;
  format: Formatters;
};

export type Value<TValue> = TValue;

export type Fn<TValue, TData> = (ctx: RecapContext<TData>) => TValue;

export type ValueOrFn<TValue, TData> = Value<TValue> | Fn<TValue, TData>;
