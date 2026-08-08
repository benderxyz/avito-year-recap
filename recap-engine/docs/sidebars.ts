import type { SidebarsConfig } from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'category',
      label: 'Старт',
      collapsed: false,
      items: ['intro', 'getting-started/quickstart', 'getting-started/architecture'],
    },
    {
      type: 'category',
      label: '@recap-engine/core',
      items: [
        'core/recap-payload',
        'core/prepare-recap',
        'core/scenes',
        'core/scene-types',
        'core/formatting',
        'core/theme',
        'core/player',
        'core/motion',
      ],
    },
    {
      type: 'category',
      label: '@recap-engine/react',
      items: [
        'react/recap',
        'react/events',
        'react/custom-scenes',
        'react/styling',
        'react/primitives',
      ],
    },
    {
      type: 'category',
      label: 'Практика',
      items: ['guides/backend-integration', 'guides/analytics-and-actions', 'guides/avito-example'],
    },
    {
      type: 'category',
      label: 'Справочник',
      collapsed: true,
      items: ['reference/core-api', 'reference/react-api'],
    },
  ],
};

export default sidebars;
