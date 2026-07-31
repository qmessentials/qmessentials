import {apiUrl} from "../api.ts";
import {Subscription} from "./types.ts";

export const subscriptionKeys = {
    root: ['subscriptions'] as const,
    all: (activeOnly: boolean = true) => [...subscriptionKeys.root, {activeOnly}] as const,
    subscription: (id: number) => [...subscriptionKeys.root, id] as const,
}

async function getSubscription(id: number): Promise<Subscription> {
    const response = await fetch(`${apiUrl}/subscription/subscriptions/${id}`, {credentials: 'include'});
    if (!response.ok) {
        throw new Error('Failed to fetch subscription');
    }
    return response.json();
}

async function getSubscriptions(activeOnly: boolean = true): Promise<Subscription[]> {
    const searchParams = new URLSearchParams({activeOnly: activeOnly.toString()});
    const response = await fetch(`${apiUrl}/subscription/subscriptions?${searchParams}`, {credentials: 'include'});
    if (!response.ok) {
        throw new Error('Failed to fetch subscriptions');
    }
    return response.json();
}

export const subscriptionQueries = {
    all: (activeOnly: boolean = true) => ({
        queryKey: subscriptionKeys.all(activeOnly),
        queryFn: () => getSubscriptions(activeOnly),
    }),
    subscription: (id: number) => ({
        queryKey: subscriptionKeys.subscription(id),
        queryFn: () => getSubscription(id),
    }),
}
