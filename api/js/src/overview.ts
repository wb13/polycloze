/* Defines elements found in overview page (e.g. language stats). */

import { createButton } from "./button";
import { getL1, setL2 } from "./data";
import { Language, LanguageStats } from "./select";

function createButtonGroup(onClick: (event: Event) => void): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("button-group");
    p.appendChild(createButton("Start", onClick));
    return p;
}

function createSeenParagraph(stats: LanguageStats): HTMLParagraphElement | null {
    const { seen, total } = stats;
    if (seen == null || seen < 0) {
        return null;
    }
    if (total == null || total <= 0) {
        return null;
    }

    const percentage = Math.floor(100 * seen / total);
    const text = `Seen: ${seen} word${seen === 1 ? "" : "s"} (${percentage}%)`;

    const p = document.createElement("p");
    p.textContent = text;
    return p;
}

function createTodayParagraph(stats: LanguageStats): HTMLParagraphElement | null {
    const { learned, reviewed } = stats;
    const texts = [];
    if (learned != null && learned > 0) {
        texts.push(`${learned} new word${learned === 1 ? "" : "s"}`);
    }
    if (reviewed != null && reviewed > 0) {
        texts.push(`${reviewed} reviewed`);
    }

    if (texts.length === 0) {
        return null;
    }

    const p = document.createElement("p");
    p.textContent = `Today: ${texts.join(", ")}`;
    return p;
}

function createLanguageStats(language: Language): HTMLDivElement | null {
    if (!language.stats) {
        return null;
    }

    const children = [];
    const seenParagraph = createSeenParagraph(language.stats);
    const todayParagraph = createTodayParagraph(language.stats);
    if (seenParagraph != null) {
        children.push(seenParagraph);
    }
    if (todayParagraph != null) {
        children.push(todayParagraph);
    }

    if (children.length === 0) {
        return null;
    }

    const div = document.createElement("div");
    div.append(...children);
    return div;
}

function createLanguageOverview(language: Language, target: string): HTMLDivElement {
    const card = document.createElement("div");
    card.classList.add("card");
    card.innerHTML = `<h2>${language.name}</h2>`;

    const row = document.createElement("div");
    row.classList.add("row");
    card.appendChild(row);

    const stats = createLanguageStats(language);
    if (stats != null) {
        row.appendChild(stats);
    }

    const start = () => {
        setL2(language.code);
        window.location.pathname = target;
    };
    row.appendChild(createButtonGroup(start));
    return card;
}

export function createOverview(languages: Language[], target = "/study"): HTMLDivElement {
    const div = document.createElement("div");
    for (const language of languages) {
        if (language.code !== getL1()) {
            div.appendChild(createLanguageOverview(language, target));
        }
    }
    return div;
}
