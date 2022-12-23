// Buttons for entering letters with diacritics.

import "./diacritic.css";
import { createButton } from "./button";
import { getL2 } from "./language";

// Enables diacritic buttons.
export function enableDiacriticButtons() {
  const lang = getL2();

  // Default value is falsy if not set, so we check the `disabled` field.
  localStorage.setItem(`diacritic.${lang.code}.disabled`, "false");
}

// Hides diacritic buttons.
export function disableDiacriticButtons() {
  const lang = getL2();
  localStorage.setItem(`diacritic.${lang.code}.disabled`, "true");
}

// Checks if diacritic buttons are enabled.
function areEnabledDiacriticButtons(): boolean {
  const lang = getL2();
  return localStorage.getItem(`diacritic.${lang.code}.disabled`) === "true"
    ? false
    : true;
}

// Creates a button that allows user to input characters with diacritics.
// Returns a button element.
function createDiacriticButton(
  name: string,
  onClick: (name: string) => void
): HTMLButtonElement {
  // TODO show equivalent digraph key presses in title tooltip
  // TODO toggle between uppercase and lowercase when shift or caps lock is
  // pressed.
  const button = createButton(name, () => onClick(name));
  button.classList.add("diacritic-button");
  button.classList.add("button-tight");
  return button;
}

// Returns array of characters to create buttons for.
function lettersWithDiacritics(languageCode: string): string[] {
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

// Returns group of diacritic buttons for the given language, or `undefined` if
// the language is not supported.
// Also returns `undefined` if diacritic buttons are disabled.
export function createDiacriticButtonGroup(
  languageCode: string,
  onClick: (name: string) => void
): HTMLParagraphElement | undefined {
  if (!areEnabledDiacriticButtons()) {
    return undefined;
  }

  const chars = lettersWithDiacritics(languageCode);
  if (chars.length <= 0) {
    return undefined;
  }

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.classList.add("diacritic-button-group");
  p.style.justifyContent = "flex-start";
  for (const char of chars) {
    p.appendChild(createDiacriticButton(char, onClick));
  }
  return p;
}
