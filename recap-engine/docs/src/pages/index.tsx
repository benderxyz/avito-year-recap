import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Heading from '@theme/Heading';
import Layout from '@theme/Layout';
import type { ReactNode } from 'react';

import styles from './index.module.css';

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={styles.heroBanner}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link className="button button--primary button--lg" to="/getting-started/quickstart">
            Быстрый старт
          </Link>
          <Link className="button button--secondary button--lg" to="/guides/backend-integration">
            Интеграция с backend
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <HomepageHeader />
      <main className={styles.features}>
        <section className="container">
          <Heading as="h2">Из данных — в интерактивную историю</Heading>
          <p>
            Типизированные сцены, готовые React layouts, темы, анимации, жесты и autoplay в двух
            небольших пакетах.
          </p>
        </section>
      </main>
    </Layout>
  );
}
