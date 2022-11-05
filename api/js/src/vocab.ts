import "./vocab.css";
import { createDateTime } from "./datetime";
import { getL2 } from "./language";
import { VocabularyItem } from "./schema";
import { createTable, createTableData, createTableHeader } from "./table";

function createVocabularyListHeader(): HTMLHeadingElement {
    const h1 = document.createElement("h1");
    const l2 = getL2();
    h1.textContent = `${l2.name} vocabulary`;
    return h1;
}

function createVocabularyListTableRow(item: VocabularyItem): HTMLTableRowElement {
    const reviewed = new Date(Date.parse(item.reviewed));
    const due = new Date(Date.parse(item.due));

    const tr = document.createElement("tr");
    tr.append(
        createTableData(item.word),
        createTableData(createDateTime(reviewed)),
        createTableData(createDateTime(due)),
        createTableData(String(item.strength)),
    );
    return tr;
}

function createVocabularyListTableBody(items: VocabularyItem[]): HTMLTableSectionElement {
    const tbody = document.createElement("tbody");
    tbody.append(...items.map(createVocabularyListTableRow));
    return tbody;
}

function createVocabularyListBody(items: VocabularyItem[]): HTMLTableElement {
    const headers = ["Word", "Last seen", "Due", "Strength"];
    return createTable(
        createTableHeader(headers),
        createVocabularyListTableBody(items),
    );
}

export function createVocabularyList(items: VocabularyItem[]): HTMLDivElement {
    const div = document.createElement("div");
    div.appendChild(createVocabularyListHeader());
    div.appendChild(createVocabularyListBody(items));
    return div;
}
