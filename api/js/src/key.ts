// Special keys.

import { createButton } from "./button";

// Creates a button that allows user to input special characters.
// Returns a button element.
function createSpecialKey(
  name: string,
  onClick: (name: string) => void
): HTMLButtonElement {
  // TODO show equivalent digraph key presses in title tooltip
  // TODO toggle between uppercase and lowercase when shift or caps lock is
  // pressed.
  const button = createButton(name, () => onClick(name));
  button.classList.add("button-tight");
  return button;
}

// Returns array of characters to create special keys for.
function specialCharacters(languageCode: string): string[] {
  switch (languageCode) {
    case "deu":
      return ["Ä", "Ö", "Ü", "ß", "ä", "é", "ö", "ü"];
    case "epo":
      return ["Ĉ", "ĉ", "Ĝ", "ĝ", "Ĥ", "ĥ", "Ĵ", "ĵ", "Ŝ", "ŝ", "Ŭ", "ŭ"];
    case "spa":
      return [
        "Á",
        "É",
        "Í",
        "Ñ",
        "Ó",
        "Ú",
        "Ü",
        "á",
        "é",
        "í",
        "ñ",
        "ó",
        "ú",
        "ü",
      ];
    default:
      return [];
  }
}

// Returns button group of special keys for given language, or `undefined` if
// the language is not supported.
export function createSpecialKeys(
  languageCode: string,
  onClick: (name: string) => void
): HTMLParagraphElement | undefined {
  const chars = specialCharacters(languageCode);
  if (chars.length <= 0) {
    return undefined;
  }

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.classList.add("button-group-dont-stretch");
  p.classList.add("special-keys");
  for (const char of chars) {
    p.appendChild(createSpecialKey(char, onClick));
  }
  return p;
}
