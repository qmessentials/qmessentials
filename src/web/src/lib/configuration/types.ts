export interface Product {
    id: number
    partNumber: string
    productName: string
    productTestConfigurations: ProductTestConfiguration[]
}

export interface Test {
    id: number
    testName: string
    documentationReferences: string[]
}

export interface ProductTestConfiguration {
    product: Product
    test: Test
    productTestSequence: number
    specificModifiers: string[]
    unit: string
    decimalPlaces: number
    minValue: number
    maxValue: number
    isCritical: boolean
}