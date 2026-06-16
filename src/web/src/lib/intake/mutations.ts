import {apiUrl} from "../api.ts";
import {SubmitTestResultInput} from "./types.ts";

async function submitTestResult(input: SubmitTestResultInput): Promise<void> {
    const response = await fetch(`${apiUrl}/intake/test-results`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        body: JSON.stringify(input),
    });
    if (!response.ok) {
        throw new Error('Failed to submit test result');
    }
}

export const testResultMutations = {
    submit: () => ({
        mutationFn: submitTestResult,
    }),
}
