/* Score counter. */

// Returns a score counter in a div.
export function createScoreCounter(correct: number): HTMLDivElement {
    const div = document.createElement("div");
    div.classList.add("medium");
    div.innerHTML = `<span class="correct">${correct}</span>`;

    const listener = (event: Event) => {
        window.removeEventListener("polycloze-count", listener);
        if ((event as CustomEvent).detail.correct) {
            div.replaceWith(createScoreCounter(correct + 1));
        } else {
            div.replaceWith(createScoreCounter(correct));
        }
    };

    window.addEventListener("polycloze-count", listener);
    return div;
}

export function dispatchUpdateCount(correct: boolean) {
    const event = new CustomEvent("polycloze-count", {
        detail: { correct }
    });
    window.dispatchEvent(event);
}
