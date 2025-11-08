export function fetchBackend(
  url: string,
  init?: RequestInit
): Promise<Response> {
  const backendUrl =
    import.meta.env.VITE_BACKEND_URL || "https://localhost:8081";
  if (!init) {
    init = {
      credentials: "include",
    };
  }
  return fetch(`${backendUrl}${url}`, init);
}

export function formatDate(date: Date): string {
  // December 31, 2023 at 5:00 PM
  return date.toLocaleString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  });
}
