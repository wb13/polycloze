import "./blank.css";
import { substituteDigraphs } from "./digraph";
import { getFont, getWidth } from "./font";

import { distance } from "fastest-levenshtein";

type Status = "correct" | "incorrect" | "almost";

function changeStatus(input: HTMLInputElement, status: Status) {
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

// modify: Invoked the first time input gets modified.
// Also returns a resize function, which should be called when the element is
// connected to the DOM.
export function createBlank(answer: string, autocapitalize: boolean, done: (answer: string, correct: boolean) => void, modify: () => void): [HTMLInputElement, () => void] {
    let correct = true;
    const input = document.createElement("input");
    input.autocapitalize = autocapitalize ? "on" : "none";
    input.classList.add("blank");

    let modified = false;
    input.addEventListener("input", () => {
        if (!modified) {
            modify();
            modified = true;
        }
        input.value = substituteDigraphs(input.value);
    });
    input.addEventListener("change", () => {
        switch (distance(input.value.trim(), answer.trim())) {
        case 0:
            changeStatus(input, "correct");
            return done(answer, correct);

        case 1:
        case 2:
            changeStatus(input, "almost");
            break;

        default:
            correct = false;
            input.placeholder = answer;
            input.value = "";
            changeStatus(input, "incorrect");
            break;
        }
    });
    return [input, () => resizeInput(input, answer)];
}
