import "./item.css";
import { createButton } from "./button";
import { Sentence, createSentence } from "./sentence";

export type Translation = {
    tatoebaID?: number
    text: string
}

export type Item = {
  sentence: Sentence
  translation: Translation
}

function showTranslationLink(translation: Translation, body: HTMLDivElement) {
    if (translation.tatoebaID == null || translation.tatoebaID <= 0) {
        return;
    }

    const a = document.createElement("a");
    a.href = "#";
    a.textContent = `#${translation.tatoebaID}`;

    const p = body.querySelector("p.translation");
    if (p != null) {
        p.textContent += " ";
        p.appendChild(a);
    }
}

function createTranslation(translation: Translation): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("translation");
    p.textContent = translation.text;
    return p;
}

function createItemBody(item: Item, done: () => void, enable: (ok: boolean) => void, clearBuffer: (frequencyClass: number) => void): [HTMLDivElement, () => void, () => void] {
    const div = document.createElement("div");
    const [sentence, check, resize] = createSentence(item.sentence, done, enable, clearBuffer);
    div.append(
        sentence,
        createTranslation(item.translation)
    );
    return [div, check, resize];
}

function createItemFooter(submitBtn: HTMLButtonElement): HTMLDivElement {
    const div = document.createElement("div");
    div.classList.add("button-group");
    div.appendChild(submitBtn);
    return div;
}

function createSubmitButton(onClick?: (event: Event) => void): [HTMLButtonElement, (ok: boolean) => void] {
    const button = createButton("Check", onClick);
    button.disabled = true;

    const enable = (ok = true) => {
        button.disabled = !ok;
    };
    return [button, enable];
}

export function createItem(item: Item, next: () => void, clearBuffer: (frequencyClass: number) => void): [HTMLDivElement, () => void] {
    const [submitBtn, enable] = createSubmitButton();

    const done = () => {
        showTranslationLink(item.translation, getBody());
        const btn = createButton("Next", next);
        submitBtn.replaceWith(btn);
        btn.focus();
    };
    const [body, check, resize] = createItemBody(item, done, enable, clearBuffer);
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
    div.append(createTranslation({text}));
    return div;
}
