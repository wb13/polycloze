import "./blank.css";
import { substituteDigraphs } from "./digraph";
import { getFont, getWidth } from "./font";
import { getL2 } from "./language";

import { distance } from "fastest-levenshtein";

type Status = "correct" | "incorrect" | "almost";

export type Answer = {
  text: string;
  normalized: string;
  new: boolean;
  difficulty: number;
};

export type Part = {
  text: string;
  answers?: Answer[];
};

export type PartWithAnswers = {
  text: string;
  answers: Answer[];
};

// Check if Part has answers.
export function hasAnswers(part: Part): boolean {
  return part.answers != null && part.answers.length > 0;
}

// Throws an exception if part has no answers.
export function requireAnswers(part: Part): PartWithAnswers {
  if (!hasAnswers(part)) {
    throw new Error("part.answers should be non-empty");
  }
  return part as PartWithAnswers;
}

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
  word = word.toLowerCase();
  return word;
}

// Count number of errors in guess.
export function compare(guess: string, answer: string): number {
  return distance(normalize(guess), normalize(answer));
}

// Sets input element status to `correct`, `almost` or `incorrect`.
// May throw an exception if `part` doesn't have answers.
//
// Summary:
// - Correct if matches exactly with a possible answer
// - Almost correct if similar to preferred answer
// - Incorrect otherwise
export function evaluateInput(
  input: HTMLInputElement,
  part: PartWithAnswers
): Status {
  const answers = part.answers;
  const diffs = [];
  for (const answer of answers) {
    const diff = compare(input.value, answer.text);
    if (diff === 0) {
      // Return immediately if exact match is found.
      changeStatus(input, "correct");
      return "correct";
    }
    diffs.push(diff);
  }

  const lang = getL2();

  // Only allow typos in preferred answer.
  if (diffs[0] <= 2 && lang.code != "jpn" && lang.code != "cmn") {
    changeStatus(input, "almost");
    return "almost";
  }

  input.placeholder = answers[0].text;
  input.value = "";
  changeStatus(input, "incorrect");
  return "incorrect";
}

// Checks if text is capitalized.
function isCapitalized(text: string): boolean {
  text = text.trim();
  const capitalized = text.toLocaleUpperCase();
  return text.charAt(0) === capitalized.charAt(0);
}

// Also returns a resize function, which should be called when the element is
// connected to the DOM.
// May throw an exception if `part` doesn't have answers.
export function createBlank(
  part: PartWithAnswers
): [HTMLInputElement, () => void] {
  const text = part.answers[0].text;

  const input = document.createElement("input");
  input.autocapitalize = isCapitalized(text) ? "on" : "none";
  input.ariaLabel = "Blank";
  input.classList.add("blank");

  input.addEventListener("input", () => {
    input.value = substituteDigraphs(input.value);
  });
  return [input, () => resizeInput(input, text)];
}
