export const apiUrl = new URL(import.meta.env.VITE_API_URL)

export async function getSession (): Promise<{ userId: string; roles: string[] } | null> {
    const response = await fetch(`${apiUrl}/session`, {credentials: 'include'})
    if (response.status === 401) return null
    if (!response.ok) throw new Error('Session request failed')
    return response.json()
}

export async function logOut () {
    const response = await fetch(`${apiUrl}/logout`, {method: 'POST', credentials: 'include'})
    if (!response.ok) throw new Error('Logout failed')
}

export async function logIn (credentials: { userId: string; password: string }) {

    const response = await fetch(`${apiUrl}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
        credentials: 'include',
    })

    if (!response.ok) {
        throw new Error('Login failed')
    }

    return response.json()
}
