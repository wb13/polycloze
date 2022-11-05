/* Course table. */

import "./course.css";
import { setL2 } from "./language";
import { Course } from "./schema";

function createCourseTableHeaderCell(content: string): HTMLTableCellElement {
    const th = document.createElement("th");
    th.scope = "col";
    th.textContent = content;
    return th;
}

function createCourseTableHeader(): HTMLTableSectionElement {
    const headers = [
        "Code",
        "Language",
        "Seen (words)",
        "Studied today",
    ];

    const tr = document.createElement("tr");
    tr.append(...headers.map(createCourseTableHeaderCell));

    const thead = document.createElement("thead");
    thead.appendChild(tr);
    return thead;
}

function createCourseTableData(child: string | Element): HTMLTableCellElement {
    const td = document.createElement("td");
    td.append(child);
    return td;
}

function createCourseMeter(course: Course): HTMLDivElement {
    const seen = course.stats?.seen || 0;
    const total = course.stats?.total || 1;

    const meter = document.createElement("meter");
    meter.max = total;
    meter.optimum = total;
    meter.value = seen;

    const description = document.createElement("div");
    description.textContent = `${seen}/${total}`;

    const div = document.createElement("div");
    div.append(meter, description);
    return div;
}

function createDiv(child: string | Element): HTMLDivElement {
    const div = document.createElement("div");
    div.append(child);
    return div;
}

function createStudiedTodayCell(course: Course): HTMLTableCellElement {
    const children = [];
    if (course.stats) {
        const { learned, reviewed } = course.stats;
        if (learned && learned > 0) {
            children.push(createDiv(`${learned} learned`));
        }
        if (reviewed && reviewed > 0) {
            children.push(createDiv(`${reviewed} reviewed`));
        }
    }

    if (children.length === 0) {
        return createCourseTableData("");
    }

    const div = document.createElement("div");
    div.style.textAlign = "center";
    div.append(...children);
    return createCourseTableData(div);
}

// Assumes course contains stats.
function createCourseTableRow(course: Course): HTMLTableRowElement {
    const tr = document.createElement("tr");
    tr.append(
        createCourseTableData(course.l2.code),
        createCourseTableData(course.l2.name),
        createCourseMeter(course),
        createStudiedTodayCell(course),
    );

    tr.addEventListener("click", () => {
        setL2(course.l2);
        window.location.pathname = "/study";
    });
    return tr;
}

function createCourseTableBody(courses: Course[]): HTMLTableSectionElement {
    const tbody = document.createElement("tbody");
    tbody.append(...courses.map(createCourseTableRow));
    return tbody;
}

export function createCourseTable(courses: Course[]): HTMLTableElement {
    const table = document.createElement("table");
    table.append(createCourseTableHeader(), createCourseTableBody(courses));
    return table;
}
