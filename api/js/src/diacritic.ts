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
  recipe?: string;
};

// Creates a button that allows user to input characters with diacritics.
// Returns a button element.
function createDiacriticButton(
  letter: Letter,
  onClick: (letter: string) => void
): HTMLButtonElement {
  const { lowercase, uppercase } = letter;

  let capsLock = false;
  let shift = false;

  const button = createButton(lowercase, () => onClick(currentValue()));
  button.classList.add("diacritic-button");
  button.classList.add("button-tight");
  if (letter.recipe != null) {
    button.title = letter.recipe;
  }

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
    case "cat":
      return [
        { uppercase: "À", lowercase: "à", recipe: "\\A` or \\a`" },
        { uppercase: "Ç", lowercase: "ç", recipe: "\\C, or \\c," },
        { uppercase: "È", lowercase: "è", recipe: "\\E` or \\e`" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Í", lowercase: "í", recipe: "\\I' or \\i'" },
        { uppercase: "Ï", lowercase: "ï", recipe: "\\I: or \\i:" },
        { uppercase: "Ò", lowercase: "ò", recipe: "\\O` or \\o`" },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o'" },
        { uppercase: "Ú", lowercase: "ú", recipe: "\\U' or \\u'" },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u:" },
      ];
    case "dan":
      return [
        { uppercase: "Å", lowercase: "å", recipe: "\\AA or \\aa" },
        { uppercase: "Æ", lowercase: "æ", recipe: "\\AE or \\ae" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ø", lowercase: "ø", recipe: "\\O/ or \\o/" },
      ];
    case "deu":
      return [
        { uppercase: "À", lowercase: "à", recipe: "\\A` or \\a`" },
        { uppercase: "Ä", lowercase: "ä", recipe: "\\A: or \\a:" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ö", lowercase: "ö", recipe: "\\O: or \\o:" },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u:" },
        { uppercase: "ß", lowercase: "ß", recipe: "\\ss" },
      ];
    case "epo":
      return [
        {
          uppercase: "Ĉ",
          lowercase: "ĉ",
          recipe: "\\C> or \\Cx or \\c> or \\cx",
        },
        {
          uppercase: "Ĝ",
          lowercase: "ĝ",
          recipe: "\\G> or \\Gx or \\g> or \\gx",
        },
        {
          uppercase: "Ĥ",
          lowercase: "ĥ",
          recipe: "\\H> or \\Hx or \\h> or \\hx",
        },
        {
          uppercase: "Ĵ",
          lowercase: "ĵ",
          recipe: "\\J> or \\Jx or \\j> or \\jx",
        },
        {
          uppercase: "Ŝ",
          lowercase: "ŝ",
          recipe: "\\S> or \\Sx or \\s> or \\sx",
        },
        {
          uppercase: "Ŭ",
          lowercase: "ŭ",
          recipe: "\\U( or \\Ux or \\u( or \\ux",
        },
      ];
    case "fin":
      return [
        { uppercase: "Å", lowercase: "å", recipe: "\\AA or \\aa" },
        { uppercase: "Ä", lowercase: "ä", recipe: "\\A: or \\a:" },
        { uppercase: "Ö", lowercase: "ö", recipe: "\\O: or \\o:" },
        { uppercase: "Š", lowercase: "š", recipe: "\\S< or \\s<" },
        { uppercase: "Ž", lowercase: "ž", recipe: "\\Z< or \\z<" },
      ];
    case "fra":
      return [
        { uppercase: "À", lowercase: "à", recipe: "\\A` or \\a`" },
        { uppercase: "Â", lowercase: "â", recipe: "\\A> or \\a>" },
        { uppercase: "Æ", lowercase: "æ", recipe: "\\AE or \\ae" },
        { uppercase: "Ç", lowercase: "ç", recipe: "\\C, or \\c," },
        { uppercase: "È", lowercase: "è", recipe: "\\E` or \\e`" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ê", lowercase: "ê", recipe: "\\E> or \\e>" },
        { uppercase: "Ë", lowercase: "ë", recipe: "\\E: or \\e:" },
        { uppercase: "Î", lowercase: "î", recipe: "\\I> or \\i>" },
        { uppercase: "Ï", lowercase: "ï", recipe: "\\I: or \\i:" },
        { uppercase: "Ô", lowercase: "ô", recipe: "\\O> or \\o>" },
        { uppercase: "Ù", lowercase: "ù", recipe: "\\U` or \\u`" },
        { uppercase: "Û", lowercase: "û", recipe: "\\U> or \\u>" },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u:" },
        { uppercase: "Ÿ", lowercase: "ÿ", recipe: "\\Y: or \\y:" },
        { uppercase: "Œ", lowercase: "œ", recipe: "\\OE or \\oe" },
      ];
    case "hrv":
      return [
        { uppercase: "Ć", lowercase: "ć", recipe: "\\C' or \\c'" },
        { uppercase: "Č", lowercase: "č", recipe: "\\C< or \\c<" },
        { uppercase: "Đ", lowercase: "đ", recipe: "\\D/ or \\d/" },
        { uppercase: "Š", lowercase: "š", recipe: "\\S< or \\s<" },
        { uppercase: "Ž", lowercase: "ž", recipe: "\\Z< or \\z<" },
      ];
    case "ita":
      return [
        { uppercase: "À", lowercase: "à", recipe: "\\A` or \\a`" },
        { uppercase: "È", lowercase: "è", recipe: "\\E` or \\e`" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ì", lowercase: "ì", recipe: "\\I` or \\i`" },
        { uppercase: "Í", lowercase: "í", recipe: "\\I' or \\i'" },
        { uppercase: "Î", lowercase: "î", recipe: "\\I> or \\i>" },
        { uppercase: "Ò", lowercase: "ò", recipe: "\\O` or \\o`" },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o'" },
        { uppercase: "Ù", lowercase: "ù", recipe: "\\U` or \\u`" },
        { uppercase: "Ú", lowercase: "ú", recipe: "\\U' or \\u'" },
      ];
    case "lit":
      return [
        { uppercase: "Ą", lowercase: "ą", recipe: "\\A; or \\a;" },
        { uppercase: "Č", lowercase: "č", recipe: "\\C< or \\c<" },
        { uppercase: "Ę", lowercase: "ę", recipe: "\\E; or \\e;" },
        { uppercase: "Ė", lowercase: "ė", recipe: "\\E. or \\e." },
        { uppercase: "Į", lowercase: "į", recipe: "\\I; or \\i;" },
        { uppercase: "Š", lowercase: "š", recipe: "\\S< or \\s<" },
        { uppercase: "Ų", lowercase: "ų", recipe: "\\U; or \\u;" },
        { uppercase: "Ū", lowercase: "ū", recipe: "\\U- or \\u-" },
        { uppercase: "Ž", lowercase: "ž", recipe: "\\Z< or \\z<" },
      ];
    case "nld":
      return [
        { uppercase: "Ĳ", lowercase: "ĳ", recipe: "\\IJ or \\ij " },
        { uppercase: "Ä", lowercase: "ä", recipe: "\\A: or \\a: " },
        { uppercase: "Ë", lowercase: "ë", recipe: "\\E: or \\e: " },
        { uppercase: "Ï", lowercase: "ï", recipe: "\\I: or \\i: " },
        { uppercase: "Ö", lowercase: "ö", recipe: "\\O: or \\o: " },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u: " },
        { uppercase: "Á", lowercase: "á", recipe: "\\A' or \\a' " },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e' " },
        { uppercase: "Í", lowercase: "í", recipe: "\\I' or \\i' " },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o' " },
        { uppercase: "Ú", lowercase: "ú", recipe: "\\U' or \\u' " },
      ];
    case "nob":
      return [
        { uppercase: "Å", lowercase: "å", recipe: "\\AA or \\aa" },
        { uppercase: "Æ", lowercase: "æ", recipe: "\\AE or \\ae" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ô", lowercase: "ô", recipe: "\\O> or \\o>" },
        { uppercase: "Ø", lowercase: "ø", recipe: "\\O/ or \\o/" },
      ];
    case "pol":
      return [
        { uppercase: "Ą", lowercase: "ą", recipe: "\\A; or \\a;" },
        { uppercase: "Ć", lowercase: "ć", recipe: "\\C' or \\c'" },
        { uppercase: "Ę", lowercase: "ę", recipe: "\\E; or \\e;" },
        { uppercase: "Ł", lowercase: "ł", recipe: "\\L/ or \\l/" },
        { uppercase: "Ń", lowercase: "ń", recipe: "\\N' or \\n'" },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o'" },
        { uppercase: "Ś", lowercase: "ś", recipe: "\\S' or \\s'" },
        { uppercase: "Ź", lowercase: "ź", recipe: "\\Z' or \\z'" },
        { uppercase: "Ż", lowercase: "ż", recipe: "\\Z. or \\z." },
      ];
    case "por":
      return [
        { uppercase: "Á", lowercase: "á", recipe: "\\A' or \\a'" },
        { uppercase: "Â", lowercase: "â", recipe: "\\A> or \\a>" },
        { uppercase: "Ã", lowercase: "ã", recipe: "\\A~ or \\a~" },
        { uppercase: "À", lowercase: "à", recipe: "\\A` or \\a`" },
        { uppercase: "Ç", lowercase: "ç", recipe: "\\C, or \\c," },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ê", lowercase: "ê", recipe: "\\E> or \\e>" },
        { uppercase: "Í", lowercase: "í", recipe: "\\I' or \\i'" },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o'" },
        { uppercase: "Ô", lowercase: "ô", recipe: "\\O> or \\o>" },
        { uppercase: "Õ", lowercase: "õ", recipe: "\\O~ or \\o~" },
        { uppercase: "Ú", lowercase: "ú", recipe: "\\U' or \\u'" },
      ];
    case "ron":
      return [
        { uppercase: "Ă", lowercase: "ă", recipe: "\\A( or \\a(" },
        { uppercase: "Â", lowercase: "â", recipe: "\\A> or \\a>" },
        { uppercase: "Î", lowercase: "î", recipe: "\\I> or \\i>" },
        { uppercase: "Ș", lowercase: "ș" },
        { uppercase: "Ț", lowercase: "ț" },
        { uppercase: "Ş", lowercase: "ş", recipe: "\\S, or \\s," },
        { uppercase: "Ţ", lowercase: "ţ", recipe: "\\T, or \\t," },
      ];
    case "spa":
      return [
        { uppercase: "Á", lowercase: "á", recipe: "\\A' or \\a'" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Í", lowercase: "í", recipe: "\\I' or \\i'" },
        {
          uppercase: "Ñ",
          lowercase: "ñ",
          recipe: "\\N? or \\ N~ or \\n? or \\n~",
        },
        { uppercase: "Ó", lowercase: "ó", recipe: "\\O' or \\o'" },
        { uppercase: "Ú", lowercase: "ú", recipe: "\\U' or \\u'" },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u:" },
      ];
    case "swe":
      return [
        { uppercase: "Á", lowercase: "á", recipe: "\\A' or \\a'" },
        { uppercase: "Ä", lowercase: "ä", recipe: "\\A: or \\a:" },
        { uppercase: "Å", lowercase: "å", recipe: "\\AA or \\aa" },
        { uppercase: "È", lowercase: "è", recipe: "\\E` or \\e`" },
        { uppercase: "É", lowercase: "é", recipe: "\\E' or \\e'" },
        { uppercase: "Ö", lowercase: "ö", recipe: "\\O: or \\o:" },
        { uppercase: "Ü", lowercase: "ü", recipe: "\\U: or \\u:" },
      ];
    case "tgl":
      return [
        {
          uppercase: "Ñ",
          lowercase: "ñ",
          recipe: "\\N? or \\ N~ or \\n? or \\n~",
        },
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
