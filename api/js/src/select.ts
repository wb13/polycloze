import "./select.css";
import { Language, getL1, setL1 } from "./language";

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

    const visited = new Map();
    for (const language of languages) {
        if (visited.has(language.code)) {
            continue;
        }
        visited.set(language.code, language);
        const option = createLanguageOption(language, language.code === getL1().code);
        select.appendChild(option);
    }

    select.addEventListener("change", () => {
        setL1(visited.get(select.value) as Language);
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
    form.style.display = "inline";
    form.append(createLanguageLabel(), createLanguageSelect(languages));
    return form;
}
