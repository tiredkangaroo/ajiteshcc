export function fetchBackend(
  url: string,
  init?: RequestInit
): Promise<Response> {
  const backendUrl =
    import.meta.env.VITE_BACKEND_URL || "https://localhost:8081";
  return fetch(`${backendUrl}${url}`, init);
}
