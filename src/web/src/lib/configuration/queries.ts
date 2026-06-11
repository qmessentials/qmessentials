import {Product} from "@/lib/configuration/types.ts";
import {apiUrl} from "@/lib/api.ts";

export const productKeys = {
    product: (partNumber: string) => ['products', partNumber] as const
}

async function getProduct(partNumber: string) : Promise<Product> {
    const response = await fetch(`${apiUrl}/config/products/${partNumber}`, {credentials: 'include'});
    return response.json();
}

export const productQueries = {
    product: (partNumber: string) => ({
        queryKey: productKeys.product(partNumber),
        queryFn: () => getProduct(partNumber),
    }),
}
