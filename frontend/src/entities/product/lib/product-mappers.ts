import type { DummyJsonProduct } from '@/shared/api';

export function formatPrice(price: number): string {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(price);
}

export function mapDummyJsonProduct(product: DummyJsonProduct) {
  return {
    id: String(product.id),
    title: product.title,
    price: product.price,
    priceLabel: formatPrice(product.price),
    category: product.category,
    thumbnail: product.thumbnail,
    rating: product.rating,
    brand: product.brand,
  };
}

export function filterProductsByPrice<T extends { price: number }>(
  products: T[],
  priceMin: number | null,
  priceMax: number | null,
): T[] {
  return products.filter((product) => {
    if (priceMin !== null && product.price < priceMin) {
      return false;
    }
    if (priceMax !== null && product.price > priceMax) {
      return false;
    }
    return true;
  });
}
