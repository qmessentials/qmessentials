import {apiUrl} from "../api.ts";
import {CreateSubscriptionInput, Subscription, UpdateSubscriptionInput} from "./types.ts";

async function createSubscription(input: CreateSubscriptionInput): Promise<Subscription> {
    const response = await fetch(`${apiUrl}/subscription/subscriptions`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify(input),
    });
    if (!response.ok) {
        throw new Error('Failed to create subscription');
    }
    return response.json();
}

async function updateSubscription(input: UpdateSubscriptionInput): Promise<Subscription> {
    const response = await fetch(`${apiUrl}/subscription/subscriptions/${input.id}`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify({ruleText: input.ruleText}),
    });
    if (!response.ok) {
        throw new Error('Failed to update subscription');
    }
    return response.json();
}

async function deleteSubscription(id: number): Promise<void> {
    const response = await fetch(`${apiUrl}/subscription/subscriptions/${id}`, {
        method: 'DELETE',
        credentials: 'include',
    });
    if (!response.ok) {
        throw new Error('Failed to delete subscription');
    }
}

export const subscriptionMutations = {
    create: () => ({
        mutationFn: createSubscription,
    }),
    update: () => ({
        mutationFn: updateSubscription,
    }),
    delete: () => ({
        mutationFn: deleteSubscription,
    }),
}
