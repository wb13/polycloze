// Returns csrf token from header meta, or empty string.
export function csrf(): string {
  const meta = document.querySelector('meta[name="csrf-token"]');
  if (meta == null) {
    return "";
  }
  return (meta as HTMLMetaElement).content;
}

export function createCSRFTokenInput(): HTMLInputElement {
  const input = document.createElement("input");
  input.type = "hidden";
  input.name = "csrf-token";
  input.value = csrf();
  return input;
}
