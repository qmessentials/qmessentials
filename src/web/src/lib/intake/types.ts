export interface Sample {
    id: number
    partNumber: string
    serialNumber: string
    status: string
    createdAt: Date
    updatedAt: Date
    testResults?: TestResult[]
}

export interface SubmitTestResultInput {
    serialNumber: string
    partNumber: string
    canonicalTestName: string
    modifiers: string[]
    testResult: number
    unit: string
    decimalPlaces: number
    minValue: number | null
    maxValue: number | null
}

export interface TestResult {
    id: string
    serialNumber: string
    partNumber: string
    canonicalTestName: string
    modifiers: string[]
    testResult: number
    unit: string
    decimalPlaces: number
    minValue: number | null
    maxValue: number | null
    hashValue: string
    voidedAt: Date | null
    voidedBy: string | null
    voidedReason: string | null
    voidComment: string | null
    createdAt: Date
    updatedAt: Date
}
