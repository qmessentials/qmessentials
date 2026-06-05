import {apiUrl} from "../api.ts";
import {Sample} from "./types.ts";

export const sampleKeys  = {
    all: ['samples'] as const,
}

async function getSamples() : Promise<Sample[]> {
    const response = await fetch(`${apiUrl}/intake/samples`, {credentials: 'include'});
    return response.json();
}

export const sampleQueries = {
    all: () => ({
        queryKey: sampleKeys.all,
        queryFn: getSamples,
    })
}