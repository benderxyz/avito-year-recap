import { describe, expect, it } from 'vitest';
import { filterProductsByPrice, formatPrice, mapDummyJsonProduct } from '@/entities/product';
import type { DummyJsonProduct } from '@/shared/api';

const sampleProduct: DummyJsonProduct = {
  id: 2,
  title: 'Eyeshadow Palette with Mirror',
  description: 'Test',
  category: 'beauty',
  price: 19.99,
  discountPercentage: 0,
  rating: 4.5,
  stock: 10,
  brand: 'Glamour Beauty',
  thumbnail: 'https://example.com/image.webp',
  images: [],
};

describe('product mappers', () => {
  it('maps dummyjson product', () => {
    const mapped = mapDummyJsonProduct(sampleProduct);
    expect(mapped.id).toBe('2');
    expect(mapped.title).toBe(sampleProduct.title);
    expect(mapped.priceLabel).toContain('₽');
  });

  it('formats price in rubles', () => {
    expect(formatPrice(1999)).toMatch(/1\s?999/);
  });

  it('filters products by price range', () => {
    const products = [{ price: 10 }, { price: 50 }, { price: 100 }];
    expect(filterProductsByPrice(products, 20, 80)).toEqual([{ price: 50 }]);
  });
});
