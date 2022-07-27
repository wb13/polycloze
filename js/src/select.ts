import "./select.css";
import { getL1, setL1 } from "./data";

export type LanguageStats = {
  seen?: number
  total?: number
  learned?: number
  reviewed?: number
  correct?: number
};

export type Language = {
  code: string
  name: string

  stats?: LanguageStats;
}

function createLanguageOption(language: Language, selected = false): HTMLOptionElement {
    const option = document.createElement("option");
    option.value = language.code;
    option.selected = selected;
    option.textContent = language.name;
    return option;
}

function createLanguageSelect(languages: Language[]): HTMLSelectElement {
    const select = document.createElement("select");
    select.id = "language-select";

    for (const language of languages) {
        const option = createLanguageOption(language, language.code === getL1());
        select.appendChild(option);
    }

    select.addEventListener("change", () => {
        setL1(select.value);
        location.reload();
    });
    return select;
}

function createLanguageLabel(): HTMLLabelElement {
    const label = document.createElement("label");
    label.htmlFor = "language-select";
    label.textContent = "🌐 ";
    return label;
}

export function createLanguageForm(languages: Language[]): HTMLFormElement {
    const form = document.createElement("form");
    form.append(createLanguageLabel(), createLanguageSelect(languages));
    return form;
}
