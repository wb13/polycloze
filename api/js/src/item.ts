import "./item.css";
import { createButton } from "./button";
import { createDiacriticButtonGroup } from "./diacritic";
import { getL1, getL2 } from "./language";
import { Sentence, createSentence } from "./sentence";
import { TTS } from "./tts";

export type Translation = {
  tatoebaID?: number;
  text: string;
};

export type Item = {
  sentence: Sentence;
  translation: Translation;
};

function showTranslationLink(translation: Translation, body: HTMLDivElement) {
  if (translation.tatoebaID == null || translation.tatoebaID <= 0) {
    return;
  }

  const a = document.createElement("a");
  a.href = `https://tatoeba.org/en/sentences/show/${translation.tatoebaID}`;
  a.target = "_blank";
  a.textContent = `#${translation.tatoebaID}`;

  const p = body.querySelector("p.translation");
  if (p != null) {
    p.textContent += " ";
    p.appendChild(a);
  }
}

// Hides diacritic buttons.
function hideDiacriticButtonGroup(body: HTMLDivElement) {
  const p = body.querySelector("p.diacritic-button-group");
  if (p != null) {
    (p as HTMLParagraphElement).style.display = "none";
  }
}

function createTranslation(translation: Translation): HTMLParagraphElement {
  const p = document.createElement("p");
  p.classList.add("translation");
  p.lang = getL1().bcp47;
  p.textContent = translation.text;
  return p;
}

function createItemBody(
  item: Item,
  done: () => void,
  enable: (ok: boolean) => void
): [HTMLDivElement, () => void, () => void] {
  const div = document.createElement("div");
  const [sentence, check, resize, inputChar] = createSentence(
    item.sentence,
    done,
    enable
  );
  div.append(sentence, createTranslation(item.translation));

  const child = createDiacriticButtonGroup(getL2().code, inputChar);
  if (child != null) {
    div.appendChild(child);
  }
  return [div, check, resize];
}

function createItemFooter(submitBtn: HTMLButtonElement): HTMLDivElement {
  const div = document.createElement("div");
  div.classList.add("button-group");
  div.appendChild(submitBtn);
  return div;
}

function createSubmitButton(
  onClick?: (event: Event) => void
): [HTMLButtonElement, (ok: boolean) => void] {
  const button = createButton("Check", onClick);
  button.classList.add("button-hidden");

  const enable = (ok = true) => {
    if (ok) {
      button.classList.remove("button-hidden");
    } else {
      button.classList.add("button-hidden");
    }
  };
  return [button, enable];
}

export function createItem(
  tts: TTS,
  item: Item,
  next: () => void
): [HTMLDivElement, () => void] {
  const [submitBtn, enable] = createSubmitButton();

  const done = () => {
    const text = item.sentence.parts.map((part) => part.text).join("");
    tts.speak(text);

    hideDiacriticButtonGroup(getBody());
    showTranslationLink(item.translation, getBody());
    const btn = createButton("Next", next);
    submitBtn.replaceWith(btn);
    btn.focus();
  };
  const [body, check, resize] = createItemBody(item, done, enable);
  const footer = createItemFooter(submitBtn);

  submitBtn.addEventListener("click", check);

  const div = document.createElement("div");
  div.classList.add("item");
  div.append(body, footer);

  function getBody(): HTMLDivElement {
    return body;
  }
  return [div, resize];
}

export function createEmptyItem(): HTMLDivElement {
  const text = "You've finished all reviews for now. Check back again later.";
  const div = document.createElement("div");
  div.classList.add("item");
  div.append(createTranslation({ text }));
  return div;
}
