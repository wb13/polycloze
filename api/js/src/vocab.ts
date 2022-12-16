import "./vocab.css";
import { fetchVocabulary } from "./api";
import { createButton } from "./button";
import { createDateTime } from "./datetime";
import { getL2 } from "./language";
import { Word } from "./schema";
import {
  createScrollingTable,
  createTable,
  createTableData,
  createTableHeader,
} from "./table";
import { TTS } from "./tts";

function createStrengthMeter(strength: number): HTMLMeterElement {
  const meter = document.createElement("meter");
  meter.min = 0;
  meter.max = 10; // ~1 year interval (log2(2 * 365))
  meter.optimum = 10;
  meter.low = 3; // ~1 week interval
  meter.high = 8; // ~6 months interval
  meter.value = strength;
  return meter;
}

function createVocabularyListHeader(): HTMLHeadingElement {
  const h1 = document.createElement("h1");
  const l2 = getL2();
  h1.textContent = `${l2.name} vocabulary`;
  return h1;
}

function createVocabularyListTableRow(
  word: Word,
  tts?: TTS
): HTMLTableRowElement {
  const learned = new Date(Date.parse(word.learned));
  const reviewed = new Date(Date.parse(word.reviewed));
  const due = new Date(Date.parse(word.due));

  const td = createTableData(word.word);
  if (tts) {
    td.addEventListener("click", () => tts.speak(word.word));
  }

  const tr = document.createElement("tr");
  tr.append(
    td,
    createTableData(createStrengthMeter(word.strength)),
    createTableData(createDateTime(learned)),
    createTableData(createDateTime(reviewed)),
    createTableData(createDateTime(due))
  );
  return tr;
}

// Creates table body for displaying vocabulary list.
// Returns a table section (tbody) and an update function for adding words to the table.
function createVocabularyListTableBody(
  tts?: TTS
): [HTMLTableSectionElement, (words: Word[]) => void] {
  const tbody = document.createElement("tbody");
  const update = (words: Word[]) =>
    tbody.append(
      ...words.map((word) => createVocabularyListTableRow(word, tts))
    );
  return [tbody, update];
}

// Creates body of vocabulary list page.
// Returns a table and an update function for adding words to the table.
function createVocabularyListBody(
  tts?: TTS
): [HTMLDivElement, (words: Word[]) => void] {
  const headers = ["Word", "Strength", "Learned", "Last seen", "Due"];
  const [body, update] = createVocabularyListTableBody(tts);
  const table = createTable(createTableHeader(headers), body);
  return [createScrollingTable(table), update];
}

function createParagraph(content: string): HTMLParagraphElement {
  const p = document.createElement("p");
  p.textContent = content;
  return p;
}

// Optionally takes a TTS object.
export async function createVocabularyList(tts?: TTS): Promise<HTMLDivElement> {
  const [body, update] = createVocabularyListBody(tts);

  const div = document.createElement("div");
  div.appendChild(createVocabularyListHeader());
  div.appendChild(body);

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.style.justifyContent = "flex-start";

  const button = p.appendChild(createButton("Load more", loadMore));
  button.style.margin = "1em 0";

  div.appendChild(p);

  let after = "";

  await loadMore();

  if (after === "") {
    // If there are no words.
    body.replaceWith(createParagraph("There's nothing to see here yet."));
  }
  return div;

  async function loadMore() {
    const items = await fetchVocabulary({ after, limit: 100 });
    const ok = items.length > 0 && items[items.length - 1].word !== after;
    if (!ok) {
      button.remove();
    } else {
      update(items);
      after = items[items.length - 1].word;
    }
  }
}
