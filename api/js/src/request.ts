// Contains functions for getting data from the server.

import { csrf } from "./csrf";

const src = findServer();

export async function fetchJson<T>(
  url: string | URL,
  options: RequestInit
): Promise<T> {
  if (url instanceof URL) {
    url = url.href;
  }
  const request = new Request(url, options);
  const response = await fetch(request);
  return await response.json();
}

export function submitJson<T>(url: string | URL, data: unknown): Promise<T> {
  const options = {
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
      "X-CSRF-Token": csrf(),
    },
    method: "POST",
    mode: "cors" as RequestMode,
  };
  return fetchJson<T>(url, options);
}

// Similar to `submitJson`, but submits a form data.
// Also returns a JSON response.
// Does not include csrf token into the form data.
// The caller should do it instead.
export async function submitFormData<T>(
  url: string | URL,
  formData: FormData
): Promise<T> {
  if (url instanceof URL) {
    url = url.href;
  }
  const options = {
    body: formData,
    method: "POST",
    mode: "cors" as RequestMode,
  };
  const request = new Request(url, options);
  const response = await fetch(request);
  return await response.json();
}

function findServer(): string {
  const url = new URL(location.href);
  if (document.currentScript == null) {
    return url.origin;
  }

  const { origin, port } = document.currentScript.dataset;
  if (origin != null) {
    return origin;
  }
  if (port != null) {
    url.port = port;
  }
  return url.origin;
}

export function resolve(path: string): URL {
  return new URL(path, src);
}
