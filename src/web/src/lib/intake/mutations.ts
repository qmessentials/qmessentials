import {apiUrl} from "../api.ts";
import {SubmitTestResultInput} from "./types.ts";

async function requestError(response: Response): Promise<Error> {
    const body = await response.json().catch(() => ({}));
    return new Error(body.error ?? `Request failed (${response.status})`);
}

async function createSample(input: {serialNumber: string; partNumber: string}) {
    const response = await fetch(`${apiUrl}/intake/samples`, {
        method: 'POST', headers: {'Content-Type': 'application/json'}, credentials: 'include', body: JSON.stringify(input),
    });
    if (!response.ok) throw await requestError(response);
    return response.json();
}

async function changeSampleStatus(input: {serialNumber: string; status: string}) {
    const response = await fetch(`${apiUrl}/intake/samples/${encodeURIComponent(input.serialNumber)}/status`, {
        method: 'PATCH', headers: {'Content-Type': 'application/json'}, credentials: 'include', body: JSON.stringify({status: input.status}),
    });
    if (!response.ok) throw await requestError(response);
    return response.json();
}

export const sampleMutations = {
    create: () => ({mutationFn: createSample}),
    changeStatus: () => ({mutationFn: changeSampleStatus}),
};

async function submitTestResult(input: SubmitTestResultInput): Promise<void> {
    const response = await fetch(`${apiUrl}/intake/test-results`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify(input),
    });
    if (!response.ok) {
        throw await requestError(response);
    }
}

export const testResultMutations = {
    submit: () => ({
        mutationFn: submitTestResult,
    }),
}
