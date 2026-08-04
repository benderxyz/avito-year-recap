import type { StoryItem } from './storyItems';

export enum EMetricType {
  Number = 'number',
  Money = 'money',
  Percentile = 'percentile',
  Ratio = 'ratio',
  String = 'string',
  List = 'list',
}

export type MetricListItem = {
  id: string;
  label: string;
  value?: number;
  imageUrl?: string;
};

type NumberMetricValue = { type: EMetricType.Number; value: number; unit?: string };

type MoneyMetricValue = { type: EMetricType.Money; value: number; currency: string };

type PercentileMetricValue = { type: EMetricType.Percentile; value: number };

type RatioMetricValue = { type: EMetricType.Ratio; value: number };

type StringMetricValue = { type: EMetricType.String; value: string };

type ListMetricValue = { type: EMetricType.List; value: MetricListItem[] };

export type MetricValue =
  | NumberMetricValue
  | MoneyMetricValue
  | PercentileMetricValue
  | RatioMetricValue
  | StringMetricValue
  | ListMetricValue;

export type NarrativeSceneCopy = {
  title?: string;
  subtitle?: string;
  body?: string;
  comparison?: string;
};

export type RecapPayloadMeta = {
  vertical: string;
  year: number;
  locale: string;
  user: {
    id: string;
    displayName: string;
    avatarUrl?: string;
  };
  generatedAt: string;
};

export type RecapPayload = {
  schemaVersion: 1;
  meta: RecapPayloadMeta;
  metrics: Record<string, MetricValue>;
  story: StoryItem[];
  narrative?: {
    scenes?: Record<string, NarrativeSceneCopy>;
    highlights?: string[];
  };
  badges?: Array<{
    id: string;
    title: string;
    description: string;
    iconUrl?: string;
  }>;
  media?: Record<string, { url: string; alt?: string }>;
  features?: {
    shareEnabled?: boolean;
    upsellId?: string;
  };
};
