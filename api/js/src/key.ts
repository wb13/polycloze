// Special keys.

import { createButton } from "./button";

// Creates a button that allows user to input special characters.
// Returns a button element.
function createSpecialKey(
  name: string,
  input: HTMLInputElement
): HTMLButtonElement {
  // TODO dispatch custom event instead of taking input
  // TODO toggle between uppercase and lowercase when shift or caps lock is
  // pressed.
  // TODO show equivalent digraph key presses in title tooltip
  return createButton(name, () => {
    input.value += name;
  });
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

// Returns button group of special keys for given language.
export function createSpecialKeys(
  languageCode: string,
  input: HTMLInputElement
): HTMLParagraphElement {
  const p = document.createElement("p");
  p.classList.add("button-group");
  for (const char of specialCharacters(languageCode)) {
    p.appendChild(createSpecialKey(char, input));
  }
  return p;
}
