import "./table.css";

function createTableHeaderCell(content: string): HTMLTableCellElement {
  const th = document.createElement("th");
  th.scope = "col";
  th.textContent = content;
  return th;
}

export function createTableHeader(headers: string[]): HTMLTableSectionElement {
  const tr = document.createElement("tr");
  tr.append(...headers.map(createTableHeaderCell));

  const thead = document.createElement("thead");
  thead.appendChild(tr);
  return thead;
}

export function createTableData(child: string | Element): HTMLTableCellElement {
  const td = document.createElement("td");
  td.append(child);
  return td;
}

export function createTable(
  header: HTMLTableSectionElement,
  body: HTMLTableSectionElement
): HTMLTableElement {
  const table = document.createElement("table");
  table.append(header, body);
  return table;
}

// Wraps around table so that it scrolls horizontally when the content overflows.
// The result is a div rather than a table.
export function createScrollingTable(table: HTMLTableElement): HTMLDivElement {
  const div = document.createElement("div");
  div.style.overflowX = "auto";
  div.appendChild(table);
  return div;
}
