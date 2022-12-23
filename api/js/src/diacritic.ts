// Buttons for entering letters with diacritics.

import "./diacritic.css";
import { createButton } from "./button";
import { createIcon } from "./icon";
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

type Letter = {
  lowercase: string;
  uppercase: string;
};

// Creates a button that allows user to input characters with diacritics.
// Returns a button element.
function createDiacriticButton(
  letter: Letter,
  onClick: (letter: string) => void
): HTMLButtonElement {
  // TODO show equivalent digraph key presses in title tooltip
  const { lowercase, uppercase } = letter;

  let capsLock = false;
  let shift = false;

  const button = createButton(lowercase, () => onClick(currentValue()));
  button.classList.add("diacritic-button");
  button.classList.add("button-tight");

  const keydownCallback = (event: KeyboardEvent) => {
    if (!button.isConnected) {
      window.removeEventListener("keydown", keydownCallback);
      return;
    }

    const before = currentValue();
    if (event.key === "Shift") {
      shift = true;
    } else if (event.key === "CapsLock") {
      if (event instanceof FakeKeyDownEvent) {
        capsLock = !capsLock;
      } else {
        const previousState = event.getModifierState("CapsLock");
        capsLock = !previousState;
      }
    } else {
      // Can't detect if caps lock is on outside of event, so just set the
      // correct value as early as possible instead.
      capsLock = event.getModifierState("CapsLock");
    }
    const after = currentValue();
    if (before !== after) {
      button.textContent = after;
    }
  };

  const keyupCallback = (event: KeyboardEvent) => {
    if (!button.isConnected) {
      window.removeEventListener("keyup", keyupCallback);
      return;
    }

    const before = currentValue();
    if (event.key === "Shift") {
      shift = false;
    }
    const after = currentValue();
    if (before !== after) {
      button.textContent = after;
    }
  };
  window.addEventListener("keydown", keydownCallback);
  window.addEventListener("keyup", keyupCallback);
  return button;

  function currentValue(): string {
    return capsLock != shift ? uppercase : lowercase;
  }
}

// Returns array of characters to create buttons for.
function lettersWithDiacritics(languageCode: string): Letter[] {
  switch (languageCode) {
    case "deu":
      return [
        { uppercase: "Ä", lowercase: "ä" },
        { uppercase: "É", lowercase: "é" },
        { uppercase: "Ö", lowercase: "ö" },
        { uppercase: "Ü", lowercase: "ü" },
        { uppercase: "ß", lowercase: "ß" },
      ];
    case "epo":
      return [
        { uppercase: "Ĉ", lowercase: "ĉ" },
        { uppercase: "Ĝ", lowercase: "ĝ" },
        { uppercase: "Ĥ", lowercase: "ĥ" },
        { uppercase: "Ĵ", lowercase: "ĵ" },
        { uppercase: "Ŝ", lowercase: "ŝ" },
        { uppercase: "Ŭ", lowercase: "ŭ" },
      ];
    case "spa":
      return [
        { uppercase: "Á", lowercase: "á" },
        { uppercase: "É", lowercase: "é" },
        { uppercase: "Í", lowercase: "í" },
        { uppercase: "Ñ", lowercase: "ñ" },
        { uppercase: "Ó", lowercase: "ó" },
        { uppercase: "Ú", lowercase: "ú" },
        { uppercase: "Ü", lowercase: "ü" },
      ];
    default:
      return [];
  }
}

class FakeKeyDownEvent extends KeyboardEvent {
  constructor(key: string) {
    super("keydown", { key });
  }
}

// Returns an on-screen caps lock toggle button.
function createCapsLockButton(): HTMLButtonElement {
  const button = createButton(createIcon("arrow-fat-up"));
  button.addEventListener("click", () => {
    // Simulate a real caps lock key press.
    window.dispatchEvent(new FakeKeyDownEvent("CapsLock"));
  });
  return button;
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

  const letters = lettersWithDiacritics(languageCode);
  if (letters.length <= 0) {
    return undefined;
  }

  const buttons = letters.map((letter) =>
    createDiacriticButton(letter, onClick)
  );

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.classList.add("diacritic-button-group");
  p.style.justifyContent = "flex-start";
  p.append(createCapsLockButton(), ...buttons);
  return p;
}

export function createDiacriticButtonSettingsSection(): HTMLFormElement {
  const form = document.createElement("form");
  form.classList.add("signin");

  form.innerHTML = `
    <div>
      <input type="checkbox" id="enable-diacritic-buttons" name="enable-diacritic-buttons">
      <label for="enable-diacritic-buttons">Enable diacritic buttons</label>
    </div>
  `;

  const input = form.querySelector("input") as HTMLInputElement;
  if (areEnabledDiacriticButtons()) {
    input.checked = true;
  }
  input.addEventListener("click", () => {
    if (input.checked) {
      enableDiacriticButtons();
    } else {
      disableDiacriticButtons();
    }
  });
  return form;
}
