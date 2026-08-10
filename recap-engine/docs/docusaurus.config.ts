import type * as Preset from '@docusaurus/preset-classic';
import type { Config } from '@docusaurus/types';
import { themes as prismThemes } from 'prism-react-renderer';

const repositoryUrl = 'https://github.com/benderxyz/avito-year-recap';
const siteUrl = process.env.DOCS_URL ?? 'https://recap-engine.netlify.app';
const baseUrl = process.env.DOCS_BASE_URL ?? '/';

const config: Config = {
  title: 'Recap Engine',
  tagline: 'Конструктор интерактивных итогов года для React',
  url: siteUrl,
  baseUrl,
  organizationName: 'benderxyz',
  projectName: 'avito-year-recap',
  onBrokenLinks: 'throw',
  markdown: {
    mermaid: true,
  },
  themes: ['@docusaurus/theme-mermaid'],
  future: {
    v4: true,
  },
  i18n: {
    defaultLocale: 'ru',
    locales: ['ru'],
  },
  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.ts',
          editUrl: `${repositoryUrl}/edit/main/recap-engine/docs/`,
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Recap Engine',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Документация',
        },
        {
          href: repositoryUrl,
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Старт',
          items: [
            { label: 'Введение', to: '/docs/intro' },
            { label: 'Быстрый старт', to: '/docs/getting-started/quickstart' },
          ],
        },
        {
          title: 'Практика',
          items: [
            { label: 'Интеграция с backend', to: '/docs/guides/backend-integration' },
            { label: 'Пример Avito app', to: '/docs/guides/avito-example' },
          ],
        },
        {
          title: 'Проект',
          items: [{ label: 'GitHub', href: repositoryUrl }],
        },
      ],
      copyright: `© ${new Date().getFullYear()} Recap Engine`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
