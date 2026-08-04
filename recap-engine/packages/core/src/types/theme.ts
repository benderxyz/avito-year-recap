export type ThemeTokens = {
  colors?: {
    bg?: string;
    fg?: string;
    muted?: string;
    accent?: string;
    accentSoft?: string;
    surface?: string;
    callout?: string;
  };
  fonts?: {
    display?: string;
    body?: string;
  };
  radii?: {
    button?: number;
    card?: number;
  };
  space?: {
    scenePadding?: number;
  };
  assets?: {
    background?: string;
  };
};

export type ResolvedTheme = {
  colors: Required<NonNullable<ThemeTokens['colors']>>;
  fonts: Required<NonNullable<ThemeTokens['fonts']>>;
  radii: Required<NonNullable<ThemeTokens['radii']>>;
  space: Required<NonNullable<ThemeTokens['space']>>;
  assets: {
    background?: string;
  };
  cssVars: Record<string, string>;
};
