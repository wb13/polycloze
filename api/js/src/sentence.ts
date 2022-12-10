import "./sentence.css";
import { submitReview } from "./api";
import { createBlank, evaluateInput, Part } from "./blank";
import { dispatchUnbuffer } from "./buffer";
import { dispatchUpdateCount } from "./counter";
import { getL2 } from "./language";
import { edit } from "./unsaved";

export type Sentence = {
    id: number
    parts: Part[]
    tatoebaID?: number
};

function createPart(text: string): HTMLSpanElement {
    const span = document.createElement("span");
    span.textContent = text;
    return span;
}

// Check of text is the beginning of a sentence.
// This is only a heuristic.
function isBeginning(text: string): boolean {
    switch (text.trim()) {
    case "":
    case "¿":
    case "¡":
    case "„":
    case "(":
    case "\"":
    case "'":
        return true;

    default:
        return false;
    }
}

// TODO document params
// Note: takes two callback functions.
// - done: ?
// - enable: Enables submit button.
// - clearBuffer: Called when frequencyClass changes to remove stale items in buffer
//
// In addition to a div element, also returns two functions to be called by the
// caller.
// - check: ?
// - resize: ?
export function createSentence(sentence: Sentence, done: () => void, enable: (ok: boolean) => void, clearBuffer: (frequencyClass: number) => void): [HTMLDivElement, () => void, () => void] {
    const resizeFns: Array<() => void> = [];
    const div = document.createElement("div");
    div.classList.add("sentence");
    div.lang = getL2().bcp47;

    for (const [i, part] of sentence.parts.entries()) {
        if (i % 2 === 0) {
            div.appendChild(createPart(part.text));
        } else {
            const autocapitalize = (i === 1) && isBeginning(sentence.parts[0].text);
            const [blank, resize] = createBlank(part, autocapitalize);
            div.appendChild(blank);
            resizeFns.push(resize);
        }
    }

    fixPunctuationWrap(div);

    const [link, render] = createSentenceLink(sentence);
    div.prepend(link);

    const check = () => {
        // Make sure everything has been filled.
        const inputs = div.querySelectorAll(".blank");
        for (let i = 0; i < inputs.length; i++) {
            const input = inputs[i] as HTMLInputElement;
            if (input.value === "") {
                return;
            }
        }

        // Time to check.
        for (let i = 0; i < inputs.length; i++) {
            const input = inputs[i];
            evaluateInput(input as HTMLInputElement, sentence.parts[2 * i + 1]);
        }

        // Check if everything is correct.
        for (let i = 0; i < inputs.length; i++) {
            const input = inputs[i];
            if (!input.classList.contains("correct")) {
                return;
            }
        }

        // Show sentence link.
        render();

        // Upload results.
        for (let i = 0; i < inputs.length; i++) {
            const input = inputs[i] as HTMLInputElement;
            const answer = input.value;

            const correct = !input.classList.contains("incorrect");
            dispatchUpdateCount(correct);
            const save = edit();
            submitReview(answer, correct).then(result => {
                dispatchUnbuffer(answer);
                save();
                clearBuffer(result.frequencyClass);
            });
        }
        div.removeEventListener("change", check);
        done();
    };
    div.addEventListener("change", check);

    div.addEventListener("input", (event: Event) => {
        if (event.target == null) {
            return;
        }
        const input = event.target as HTMLInputElement;
        enable(input.value !== "");
    });

    const resizeAll = () => {
        for (const fn of resizeFns) {
            fn();
        }
    };
    return [div, check, resizeAll];
}

// Prevents punctuation symbols from starting a new line.
// Assumes all child nodes are elements.
function fixPunctuationWrap(div: HTMLDivElement) {
    const inputs = div.querySelectorAll(".blank");
    for (let i = 0; i < inputs.length; i++) {
        const input = inputs[i];
        const span = input.nextElementSibling;

        if (span == null) {
            continue;
        }

        // NOTE Does not split by other whitespace characters
        const words = span.textContent?.split(" ") || [];
        if (words.length > 0 && words[0] !== "") {
            const wrapper = document.createElement("span");
            wrapper.style.whiteSpace = "nowrap";
            input.replaceWith(wrapper);
            wrapper.appendChild(input);

            const after = document.createElement("span");
            after.textContent = words[0];
            wrapper.appendChild(after);

            words.shift();

            const tail = words.join(" ");
            span.textContent = tail.length > 0 ? " " + tail : "";
        }
    }
}

function createSentenceLink(sentence: Sentence): [HTMLDivElement, () => void] {
    const div = document.createElement("div");
    div.classList.add("sentence-link");
    div.classList.add("transparent");
    div.textContent = "#";

    const render = () => {
        if (sentence.tatoebaID == null || sentence.tatoebaID <= 0) {
            return;
        }
        const url = `https://tatoeba.org/en/sentences/show/${sentence.tatoebaID}`;
        div.innerHTML = `<a href="${url}" target="_blank">#${sentence.tatoebaID}</a>`;
        div.classList.remove("transparent");
    };
    return [div, render];
}
