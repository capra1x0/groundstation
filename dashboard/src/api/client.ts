export const API_URL: string = import.meta.env.VITE_API_URL ?? "http://localhost:1880";

export async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${API_URL}${path}`);

  if (!response.ok) {
    throw new Error(`GET ${path} failed: ${response.status} ${response.statusText}`);
  }

  return (await response.json()) as T;
}