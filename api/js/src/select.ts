import "./select.css";
import { getL1, setL1 } from "./data";

export type Language = {
  code: string
  name: string
  bcp47: string
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

    const visited = new Set();
    for (const language of languages) {
        if (visited.has(language.code)) {
            continue;
        }
        visited.add(language.code);
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
