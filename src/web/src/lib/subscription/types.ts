export interface Subscription {
    id: number
    ownerUserId: string
    ruleText: string
    versionId: number
    isActive: boolean
    createdAt: Date
    updatedAt: Date
}

export interface CreateSubscriptionInput {
    ruleText: string
}

export interface UpdateSubscriptionInput {
    id: number
    ruleText: string
}
