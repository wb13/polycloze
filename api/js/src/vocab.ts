import "./vocab.css";
import { fetchVocabulary } from "./api";
import { createButton } from "./button";
import { createDateTime } from "./datetime";
import { getL2 } from "./language";
import { Word } from "./schema";
import { createScrollingTable, createTable, createTableData, createTableHeader } from "./table";

function createStrengthMeter(strength: number): HTMLMeterElement {
    const meter = document.createElement("meter");
    meter.min = 0;
    meter.max = 10;  // ~1 year interval (log2(2 * 365))
    meter.optimum = 10;
    meter.low = 3;  // ~1 week interval
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

function createVocabularyListTableRow(word: Word): HTMLTableRowElement {
    const learned = new Date(Date.parse(word.learned));
    const reviewed = new Date(Date.parse(word.reviewed));
    const due = new Date(Date.parse(word.due));

    const tr = document.createElement("tr");
    tr.append(
        createTableData(word.word),
        createTableData(createDateTime(learned)),
        createTableData(createDateTime(reviewed)),
        createTableData(createDateTime(due)),
        createTableData(createStrengthMeter(word.strength)),
    );
    return tr;
}

// Creates table body for displaying vocabulary list.
// Returns a table section (tbody) and an update function for adding words to the table.
function createVocabularyListTableBody(): [HTMLTableSectionElement, (words: Word[]) => void] {
    const tbody = document.createElement("tbody");
    const update = (words: Word[]) => tbody.append(...words.map(createVocabularyListTableRow));
    return [tbody, update];
}

// Creates body of vocabulary list page.
// Returns a table and an update function for adding words to the table.
function createVocabularyListBody(): [HTMLDivElement, (words: Word[]) => void] {
    const headers = ["Word", "Learned", "Last seen", "Due", "Strength"];
    const [body, update] = createVocabularyListTableBody();
    const table = createTable(createTableHeader(headers), body);
    return [createScrollingTable(table), update];
}

function createParagraph(content: string): HTMLParagraphElement {
    const p = document.createElement("p");
    p.textContent = content;
    return p;
}

export async function createVocabularyList(): Promise<HTMLDivElement> {
    const [body, update] = createVocabularyListBody();

    const div = document.createElement("div");
    div.appendChild(createVocabularyListHeader());
    div.appendChild(body);
    const button = div.appendChild(createButton("Load more", loadMore));
    button.style.margin = "1em 0";

    let after = "";

    await loadMore();

    if (after === "") {
        // If there are no words.
        body.replaceWith(createParagraph("There's nothing to see here."));
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
