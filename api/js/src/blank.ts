import "./blank.css";
import { substituteDigraphs } from "./digraph";
import { getFont, getWidth } from "./font";

import { distance } from "fastest-levenshtein";

type Status = "correct" | "incorrect" | "almost";

export type Answer = {
    text: string;
    normalized: string;
};

export type Part = {
    text: string;
    answers?: Answer[];
};

function changeStatus(input: HTMLInputElement, status: Status) {
    // input.classList.remove("correct");
    // input.classList.remove("almost");
    // input.classList.remove("incorrect");
    input.classList.add(status);
}

// Resize input element to fit text.
function resizeInput(input: HTMLInputElement, text: string) {
    if (!input.isConnected) {
        console.error("should only be called on connected elements");
    }
    const width = getWidth(getFont(input), text);
    input.style.setProperty("width", width);
}

// Removes soft-hyphens and unnecessary surrounding characters (whitespace,
// zero-width spaces, no-break spaces, etc.
function normalize(word: string): string {
    // Remove soft-hyphens.
    word = word.trim().replace(/\u00AD/g, "");

    const zeroWidthSpace = "\u200B";
    while (word.startsWith(zeroWidthSpace)) {
        word = word.slice(zeroWidthSpace.length);
    }
    while (word.endsWith(zeroWidthSpace)) {
        word = word.slice(0, word.length - zeroWidthSpace.length);
    }

    const noBreakSpace = "\u00A0";
    while (word.startsWith(noBreakSpace)) {
        word = word.slice(noBreakSpace.length);
    }
    while (word.endsWith(noBreakSpace)) {
        word = word.slice(0, word.length - noBreakSpace.length);
    }
    return word;
}

// Count number of errors in guess.
function compare(guess: string, answer: string): number {
    return distance(normalize(guess), normalize(answer));
}

export function evaluateInput(input: HTMLInputElement, part: Part): Status {
    // TODO compare with every possible answer
    const answer = (part.answers as Answer[])[0].text;
    switch (compare(input.value, answer)) {
    case 0:
        changeStatus(input, "correct");
        return "correct";

    case 1:
    case 2:
        changeStatus(input, "almost");
        return "almost";

    default:
        input.placeholder = answer;
        input.value = "";
        changeStatus(input, "incorrect");
        return "incorrect";
    }
}

// Also returns a resize function, which should be called when the element is
// connected to the DOM.
export function createBlank(part: Part, autocapitalize: boolean): [HTMLInputElement, () => void] {
    const input = document.createElement("input");
    input.autocapitalize = autocapitalize ? "on" : "none";
    input.ariaLabel = "Blank";
    input.classList.add("blank");

    input.addEventListener("input", () => {
        input.value = substituteDigraphs(input.value);
    });

    const answer = (part.answers || [])[0];
    const text = answer?.text || "abcdef";
    // TODO change input type to PartWithAnswers instead of using placeholder text
    return [input, () => resizeInput(input, text)];
}
