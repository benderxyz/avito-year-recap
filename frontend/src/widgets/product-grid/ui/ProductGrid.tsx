import { Badge, Card, Group, Image, SimpleGrid, Stack, Text } from '@mantine/core';
import type { Product } from '@/entities/product';
import styles from './ProductGrid.module.css';

type ProductGridProps = {
  products: Product[];
  title?: string;
  subtitle?: string;
};

export function ProductGrid({
  products,
  title = 'Рекомендации для вас',
  subtitle = 'Объявления рядом и по интересам',
}: ProductGridProps) {
  return (
    <section className={styles.section}>
      <Stack gap={4} mb="md">
        <Text component="h1" className={styles.title}>
          {title}
        </Text>
        <Text className={styles.subtitle}>{subtitle}</Text>
      </Stack>

      <SimpleGrid cols={{ base: 1, xs: 2, sm: 3, md: 4 }} spacing="md">
        {products.map((product) => (
          <Card key={product.id} className={styles.card} padding={0} radius="md" withBorder>
            <Card.Section>
              <Image src={product.thumbnail} alt={product.title} height={180} fit="cover" />
            </Card.Section>
            <Stack gap={6} p="sm">
              <Text className={styles.price}>{product.priceLabel}</Text>
              <Text className={styles.productTitle} lineClamp={2}>
                {product.title}
              </Text>
              <Group justify="space-between" gap="xs">
                <Text className={styles.meta}>{product.brand}</Text>
                <Badge variant="light" color="avito">
                  ★ {product.rating.toFixed(1)}
                </Badge>
              </Group>
            </Stack>
          </Card>
        ))}
      </SimpleGrid>
    </section>
  );
}
