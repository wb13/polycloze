import "./vocab.css";
import { createDateTime } from "./datetime";
import { getL2 } from "./language";
import { VocabularyItem } from "./schema";

function createVocabularyListHeader(): HTMLHeadingElement {
    const h1 = document.createElement("h1");
    const l2 = getL2();
    h1.textContent = `${l2.name} vocabulary`;
    return h1;
}

function createVocabularyListTableHeaderCell(content: string): HTMLTableCellElement {
    const th = document.createElement("th");
    th.scope = "col";
    th.textContent = content;
    return th;
}

function createVocabularyListTableHeader(): HTMLTableSectionElement {
    const headers = [
        "Word",
        "Last seen",
        "Due",
        "Strength",
    ];

    const tr = document.createElement("tr");
    tr.append(...headers.map(createVocabularyListTableHeaderCell));

    const thead = document.createElement("thead");
    thead.appendChild(tr);
    return thead;
}

function createVocabularyListTableData(child: string | Element): HTMLTableCellElement {
    const td = document.createElement("td");
    td.append(child);
    return td;
}

function createVocabularyListTableRow(item: VocabularyItem): HTMLTableRowElement {
    const reviewed = new Date(Date.parse(item.reviewed));
    const due = new Date(Date.parse(item.due));

    const tr = document.createElement("tr");
    tr.append(
        createVocabularyListTableData(item.word),
        createVocabularyListTableData(createDateTime(reviewed)),
        createVocabularyListTableData(createDateTime(due)),
        createVocabularyListTableData(String(item.strength)),
    );
    return tr;
}

function createVocabularyListTableBody(items: VocabularyItem[]): HTMLTableSectionElement {
    const tbody = document.createElement("tbody");
    tbody.append(...items.map(createVocabularyListTableRow));
    return tbody;
}

function createVocabularyListTable(items: VocabularyItem[]): HTMLTableElement {
    const table = document.createElement("table");
    table.append(createVocabularyListTableHeader(), createVocabularyListTableBody(items));
    return table;
}

function createVocabularyListBody(items: VocabularyItem[]): HTMLTableElement {
    return createVocabularyListTable(items);
}

export function createVocabularyList(items: VocabularyItem[]): HTMLDivElement {
    const div = document.createElement("div");
    div.appendChild(createVocabularyListHeader());
    div.appendChild(createVocabularyListBody(items));
    return div;
}
