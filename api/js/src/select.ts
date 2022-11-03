import "./select.css";
import { createButton } from "./button";
import { onClickOutside } from "./click";
import { Language, getL1, setL1 } from "./language";

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
// Returns a div element and a function that toggles the visibility of the menu.
// The menu is hidden by default (display: none).
function createLanguageMenu(languages: Language[]): [HTMLDivElement, () => void] {
    const div = document.createElement("div");
    div.classList.add("menu");
    div.style.display = "none";
    const toggle = () => {
        div.style.display = div.style.display === "none" ? "" : "none";
    };

    const current = getL1();
    const visited = new Map();
    for (const language of languages) {
        if (current.code === language.code || visited.has(language.code)) {
            continue;
        }
        visited.set(language.code, language);
        div.appendChild(createLanguageButton(language));
    }
    return [div, toggle];
}

export function createLanguageSelectButton(languages: Language[]): HTMLButtonElement {
    const [menu, toggle] = createLanguageMenu(languages);
    document.body.appendChild(menu);

    const l1 = getL1();
    const content = `🌐 ${l1.name}`;
    const button = createButton(content, toggle);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");

    // Listen for clicks outside the menu.
    document.addEventListener("click", onClickOutside(menu, (target: EventTarget) => {
        if (target != button) {
            menu.style.display = "none";
        }
    }));
    return button;
}
