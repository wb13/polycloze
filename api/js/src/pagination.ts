import "./pagination.css";
import { createButton } from "./button";

function createTightButton(content: string | Element, onClick?: (event: Event) => void): HTMLButtonElement {
    const button = createButton(content, onClick);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}

function createPreviousButton(disabled: boolean, onClick: (event: Event) => void): HTMLButtonElement {
    const button = createTightButton("Prev", onClick);
    button.disabled = disabled;
    return button;
}

function createNextButton(disabled: boolean, onClick: (event: Event) => void): HTMLButtonElement {
    const button = createTightButton("Next", onClick);
    button.disabled = disabled;
    return button;
}

function createPageButton(page: number, selected: boolean, onClick: (event: Event) => void): HTMLButtonElement {
    if (!selected) {
        return createTightButton(String(page), onClick);
    }
    const content = document.createElement("b");
    content.textContent = String(page);
    return createTightButton(content);
}

export function createPagination(page: number, lastPage: number, gotoPage: (page: number) => void): HTMLDivElement {
    if (page <= 0) {
        throw new Error(`page should be positive: ${page}`);
    }
    if (lastPage <= 0) {
        throw new Error(`lastPage should be positive: ${lastPage}`);
    }
    if (page > lastPage) {
        throw new Error(`page should be <= lastPage: ${page}, ${lastPage}`);
    }
    const children: Array<string | Element> = [createPreviousButton(page === 1, () => gotoPage(page - 1))];

    children.push(createPageButton(1, page === 1, () => gotoPage(1)));
    if (page >= 5) {
        children.push("⋯");
    } else if (page === 4) {
        children.push(createPageButton(2, false, () => gotoPage(2)));
    }
    for (let i = Math.max(2, page - 1); i <= Math.min(page + 1, lastPage); i++) {
        children.push(createPageButton(i, page === i, () => gotoPage(i)));
    }
    if (lastPage - page >= 4) {
        children.push("⋯");
    } else if (lastPage - page === 3) {
        children.push(createPageButton(lastPage - 1, page === lastPage - 1, () => gotoPage(lastPage - 1)));
    }
    if (lastPage - page >= 2) {
        children.push(createPageButton(lastPage, page === lastPage, () => gotoPage(lastPage)));
    }

    children.push(createNextButton(page === lastPage, () => gotoPage(page + 1)));

    const div = document.createElement("div");
    div.classList.add("pagination");
    div.append(...children);
    return div;
}
