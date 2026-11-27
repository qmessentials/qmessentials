import {Product} from "@/lib/configuration/types.ts";
import {apiUrl} from "@/lib/api.ts";

export const productKeys = {
    all: (activeOnly: boolean) => ['products', {activeOnly}] as const,
    product: (partNumber: string) => ['products', partNumber] as const
}

async function getProducts(activeOnly: boolean): Promise<Product[]> {
    const response = await fetch(`${apiUrl}/config/products?activeOnly=${activeOnly}`, {credentials: 'include'});
    if (!response.ok) throw new Error('Could not load products');
    return response.json();
}

async function getProduct(partNumber: string) : Promise<Product> {
    const response = await fetch(`${apiUrl}/config/products/${partNumber}`, {credentials: 'include'});
    return response.json();
}

export const productQueries = {
    all: (activeOnly = true) => ({
        queryKey: productKeys.all(activeOnly),
        queryFn: () => getProducts(activeOnly),
    }),
    product: (partNumber: string) => ({
        queryKey: productKeys.product(partNumber),
        queryFn: () => getProduct(partNumber),
    }),
}
