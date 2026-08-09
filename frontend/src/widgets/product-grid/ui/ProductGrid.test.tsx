import { MantineProvider } from '@mantine/core';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import type { Product } from '@/entities/product';
import { mantineTheme } from '@/shared/ui';
import { ProductGrid } from '@/widgets/product-grid';

const products: Product[] = [
  {
    id: '1',
    title: 'Test product',
    price: 1000,
    priceLabel: '1 000 ₽',
    category: 'beauty',
    thumbnail: 'https://example.com/image.webp',
    rating: 4.8,
    brand: 'Brand',
  },
];

describe('ProductGrid', () => {
  it('renders product cards', () => {
    render(
      <MantineProvider theme={mantineTheme}>
        <ProductGrid products={products} />
      </MantineProvider>,
    );

    expect(screen.getByText('Test product')).toBeInTheDocument();
    expect(screen.getByText('1 000 ₽')).toBeInTheDocument();
  });
});
