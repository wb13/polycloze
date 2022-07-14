import "./item.css";
import { createButton } from "./button";
import { Sentence, createSentence } from "./sentence";

export type Item = {
  sentence: Sentence
  translation: string
}

function createTranslation(translation: string): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("translation");
    p.textContent = translation;
    return p;
}

function createItemBody(item: Item, next: (ok: boolean) => void, enable: () => void): [HTMLDivElement, () => void, () => void] {
    const div = document.createElement("div");
    const [sentence, check, resize] = createSentence(item.sentence, next, enable);
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

function createSubmitButton(onClick?: (event: Event) => void): [HTMLButtonElement, () => void] {
    const button = createButton("Check", onClick);
    button.disabled = true;

    const enable = () => {
        button.disabled = false;
    };
    return [button, enable];
}

export function createItem(item: Item, next: (ok: boolean) => void): [HTMLDivElement, () => void] {
    const [submitBtn, enable] = createSubmitButton();

    const [body, check, resize] = createItemBody(item, next, enable);
    const footer = createItemFooter(submitBtn);

    submitBtn.addEventListener("click", check);

    const div = document.createElement("div");
    div.classList.add("item");
    div.append(body, footer);
    return [div, resize];
}

export function createEmptyItem(): HTMLDivElement {
    const text = "You've finished all reviews for now. Check back again later.";
    const div = document.createElement("div");
    div.classList.add("item");
    div.append(createTranslation(text));
    return div;
}
