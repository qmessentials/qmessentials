export const apiUrl = new URL(import.meta.env.VITE_API_URL)

export async function logIn () {

    const response = await fetch(`${apiUrl}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({}),
        credentials: 'include',
    })

    if (!response.ok) {
        throw new Error('Login failed')
    }

    return response.json()
}
