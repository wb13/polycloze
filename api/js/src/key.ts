// Special keys.

import { createButton } from "./button";

function announceCharacterClick(value: string) {
  const event = new CustomEvent("polycloze-special-character", {
    detail: value,
  });
  window.dispatchEvent(event);
}

// Creates a button that allows user to input special characters.
// Returns a button element.
// Emits a `polycloze-special-character` event when clicked.
function createSpecialKey(name: string): HTMLButtonElement {
  // TODO show equivalent digraph key presses in title tooltip
  // TODO toggle between uppercase and lowercase when shift or caps lock is
  // pressed.
  return createButton(name, () => announceCharacterClick(name));
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
// Emits a `polycloze-special-character` custom event when one of the buttons
// is clicked.
export function createSpecialKeys(languageCode: string): HTMLParagraphElement {
  const p = document.createElement("p");
  p.classList.add("button-group");
  for (const char of specialCharacters(languageCode)) {
    p.appendChild(createSpecialKey(char));
  }
  return p;
}
