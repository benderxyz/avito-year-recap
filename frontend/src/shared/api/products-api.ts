import { DUMMYJSON_BASE_URL, httpClient } from './http';

export type DummyJsonProduct = {
  id: number;
  title: string;
  description: string;
  category: string;
  price: number;
  discountPercentage: number;
  rating: number;
  stock: number;
  brand: string;
  thumbnail: string;
  images: string[];
};

export type DummyJsonProductsResponse = {
  products: DummyJsonProduct[];
  total: number;
  skip: number;
  limit: number;
};

export type ProductsQueryParams = {
  q?: string;
  category?: string;
  sortBy?: string;
  order?: 'asc' | 'desc';
  limit?: number;
  skip?: number;
};

export async function fetchProducts(
  params: ProductsQueryParams,
): Promise<DummyJsonProductsResponse> {
  const { q, category, sortBy, order, limit = 12, skip = 0 } = params;

  if (q?.trim()) {
    const response = await httpClient.get<DummyJsonProductsResponse>(
      `${DUMMYJSON_BASE_URL}/products/search`,
      { params: { q: q.trim(), limit, skip, sortBy, order } },
    );
    return response.data;
  }

  if (category) {
    const response = await httpClient.get<DummyJsonProductsResponse>(
      `${DUMMYJSON_BASE_URL}/products/category/${category}`,
      { params: { limit, skip, sortBy, order } },
    );
    return response.data;
  }

  const response = await httpClient.get<DummyJsonProductsResponse>(
    `${DUMMYJSON_BASE_URL}/products`,
    { params: { limit, skip, sortBy, order } },
  );
  return response.data;
}

export async function fetchProductCategories(): Promise<string[]> {
  const response = await httpClient.get<string[]>(`${DUMMYJSON_BASE_URL}/products/category-list`);
  return response.data;
}
