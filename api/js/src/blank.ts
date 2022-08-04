import "./blank.css";
import { substituteDigraphs } from "./digraph";
import { getFont, getWidth } from "./font";

import { distance } from "fastest-levenshtein";

type Status = "correct" | "incorrect" | "almost";

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

export function evaluateInput(input: HTMLInputElement, answer: string): Status {
    switch (distance(input.value.trim(), answer.trim())) {
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
export function createBlank(answer: string, autocapitalize: boolean): [HTMLInputElement, () => void] {
    const input = document.createElement("input");
    input.autocapitalize = autocapitalize ? "on" : "none";
    input.ariaLabel = "Blank";
    input.classList.add("blank");

    input.addEventListener("input", () => {
        input.value = substituteDigraphs(input.value);
    });
    return [input, () => resizeInput(input, answer)];
}
