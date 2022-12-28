import "./sentence.css";
import {
  compare,
  createBlank,
  evaluateInput,
  hasAnswers,
  Part,
  PartWithAnswers,
} from "./blank";
import { announceResult } from "./buffer";
import { getL2 } from "./language";

export type Sentence = {
  id: number;
  parts: Part[];
  tatoebaID?: number;
};

function createPart(text: string): HTMLSpanElement {
  const span = document.createElement("span");
  span.textContent = text;
  return span;
}

// TODO document params
// Note: takes two callback functions.
// - done: ?
// - enable: Enables submit button.
//
// In addition to a div element, also returns functions to be called by the
// caller.
// - check: ?
// - resize: ?
// - `inputChar`: append a character to the last input element in focus.
export function createSentence(
  sentence: Sentence,
  done: () => void,
  enable: (ok: boolean) => void
): [HTMLDivElement, () => void, () => void, (char: string) => void] {
  const resizeFns: Array<() => void> = [];
  const div = document.createElement("div");
  div.classList.add("sentence");
  div.lang = getL2().bcp47;

  // Last focused blank. Diacritic buttons append to this input element if not
  // null.
  let lastFocused: HTMLInputElement | null = null;

  // NOTE `inputs` and `blankParts` have the same length.
  const inputs: HTMLInputElement[] = [];
  const blankParts: PartWithAnswers[] = [];
  for (const part of sentence.parts) {
    if (!hasAnswers(part)) {
      div.appendChild(createPart(part.text));
      continue;
    }

    const checkedPart = part as PartWithAnswers;
    const [blank, resize] = createBlank(checkedPart);
    div.appendChild(blank);
    resizeFns.push(resize);

    inputs.push(blank);
    blankParts.push(checkedPart);

    // Add event listener to input element to update `lastFocused`.
    blank.addEventListener("focus", () => {
      lastFocused = blank;
    });
  }

  fixPunctuationWrap(div);

  const [link, render] = createSentenceLink(sentence);
  div.prepend(link);

  const check = () => {
    // False-positive event if a diacritic button is active.
    // This happens because clicking on these buttons removes the focus from
    // the input element, which triggers a "change" event.
    const activeElement = document.activeElement;
    if (activeElement?.tagName === "BODY") {
      return;
    }
    if (
      activeElement != null &&
      activeElement.tagName === "BUTTON" &&
      activeElement.classList.contains("diacritic-button")
    ) {
      return;
    }

    // Make sure everything has been filled.
    if (inputs.some((input) => input.value === "")) {
      return;
    }

    // Time to check.
    for (const [i, input] of inputs.entries()) {
      evaluateInput(input, blankParts[i]);
    }

    // Check if everything is correct.
    if (inputs.some((input) => !input.classList.contains("correct"))) {
      return;
    }

    // Show sentence link.
    render();

    // Notify buffer of results.
    for (const [i, input] of inputs.entries()) {
      const answer = input.value;

      // Normalize word.
      let word = answer;
      let new_ = false;
      for (const answer of blankParts[i].answers) {
        if (compare(input.value, answer.text) === 0) {
          word = answer.normalized;
          new_ = answer.new;
          break;
        }
      }

      const correct = !input.classList.contains("incorrect");
      announceResult({
        word,
        correct,
        new: new_,
        timestamp: Math.floor(Date.now() / 1000),
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
  const inputChar = (char: string) => {
    // Replace selection with character.
    if (lastFocused != null) {
      const value = lastFocused.value;
      const start = lastFocused.selectionStart || value.length;
      const end = lastFocused.selectionEnd || value.length;
      lastFocused.value = value.slice(0, start) + char + value.slice(end);
      lastFocused.focus();
    }
  };
  return [div, check, resizeAll, inputChar];
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
