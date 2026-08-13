export const routes = {
  badges: '/badges',
  badgeCreate: '/badges/new',
  badgeById: '/badges/$id',
  metrics: '/metrics',
  metricCreate: '/metrics/new',
  metricByKey: '/metrics/$key',
  recommendations: '/recommendations',
  recommendationById: '/recommendations/$id',
  stories: '/stories',
  storyById: '/stories/$id',
  preview: '/preview',
} as const;
