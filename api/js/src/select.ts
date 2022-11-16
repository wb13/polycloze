import "./select.css";
import { createButton } from "./button";
import { createIcon } from "./icon";
import { getL1, setL1 } from "./language";
import { createModal } from "./modal";
import { Language } from "./schema";

// Creates an entry in the language select menu.
function createLanguageButton(language: Language): HTMLButtonElement {
    const button = createButton(language.name, () => {
        setL1(language);
        location.reload();
    });
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}

// Creates a language select menu/modal
function createLanguageMenu(languages: Language[]): [HTMLDivElement, () => void] {
    const body = document.createElement("div");
    body.classList.add("menu");

    const current = getL1();
    const visited = new Map();
    for (const language of languages) {
        if (current.code === language.code || visited.has(language.code)) {
            continue;
        }
        visited.set(language.code, language);
        body.appendChild(createLanguageButton(language));
    }
    return createModal(body);
}

export function createLanguageSelectButton(languages: Language[]): HTMLButtonElement {
    const [menu, show] = createLanguageMenu(languages);
    document.body.appendChild(menu);

    const l1 = getL1();
    const content = document.createElement("span");
    content.append(createIcon("globe"), ` ${l1.name}`);

    const button = createButton(content, show);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}
