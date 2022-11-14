import "./select.css";
import { createButton } from "./button";
import { onClickOutside } from "./click";
import { createIcon } from "./icon";
import { getL1, setL1 } from "./language";
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

// Creates a language select menu.
// Returns:
// - a div element
// - a function that shows the menu
// - a function that hides the menu.
// The menu is hidden by default (display: none).
function createLanguageMenu(languages: Language[]): [HTMLDivElement, () => void, () => void] {
    const div = document.createElement("div");
    div.classList.add("menu");
    div.style.display = "none";
    const show = () => {
        div.style.display = "";
    };
    const hide = () => {
        div.style.display = "none";
    };

    // Create close button.
    const button = createButton("✕", hide);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    div.appendChild(button);

    const current = getL1();
    const visited = new Map();
    for (const language of languages) {
        if (current.code === language.code || visited.has(language.code)) {
            continue;
        }
        visited.set(language.code, language);
        div.appendChild(createLanguageButton(language));
    }
    return [div, show, hide];
}

export function createLanguageSelectButton(languages: Language[]): HTMLButtonElement {
    const [menu, show, hide] = createLanguageMenu(languages);
    document.body.appendChild(menu);

    const l1 = getL1();
    const content = document.createElement("span");
    content.append(createIcon("globe"), ` ${l1.name}`);

    const button = createButton(content, show);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");

    // Listen for clicks outside the menu.
    document.addEventListener("click", onClickOutside(menu, (target: EventTarget) => {
        if (target != button) {
            hide();
        }
    }));
    return button;
}
