/* Course table. */

import "./course.css";
import { createIcon } from "./icon";
import { setL2 } from "./language";
import { Course } from "./schema";
import { createScrollingTable, createTable, createTableData, createTableHeader } from "./table";

function createSeenCountCell(course: Course): HTMLTableCellElement {
    const td = createTableData("");
    const seen = course.stats?.seen || 0;
    if (seen > 0) {
        td.append(createIcon("info"), ` ${seen} words`);
        td.addEventListener("click", () => {
            setL2(course.l2);
            window.location.pathname = "/vocab";
        });
    }
    return td;
}

function createDiv(child: string | Element): HTMLDivElement {
    const div = document.createElement("div");
    div.append(child);
    return div;
}

function createCourseCodeCell(course: Course): HTMLTableCellElement {
    const td = createTableData("");
    td.append(createIcon("play-circle"), ` ${course.l2.code}`);
    td.addEventListener("click", () => {
        setL2(course.l2);
        window.location.pathname = "/study";
    });
    return td;
}

function createCourseNameCell(course: Course): HTMLTableCellElement {
    const td = createTableData(course.l2.name);
    td.addEventListener("click", () => {
        setL2(course.l2);
        window.location.pathname = "/study";
    });
    return td;
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
        return createTableData("");
    }

    const div = document.createElement("div");
    div.style.textAlign = "center";
    div.append(...children);
    return createTableData(div);
}

// Assumes course contains stats.
function createCourseTableRow(course: Course): HTMLTableRowElement {
    const tr = document.createElement("tr");
    tr.append(
        createCourseCodeCell(course),
        createCourseNameCell(course),
        createSeenCountCell(course),
        createStudiedTodayCell(course),
    );
    return tr;
}

function createCourseTableBody(courses: Course[]): HTMLTableSectionElement {
    const tbody = document.createElement("tbody");
    tbody.append(...courses.map(createCourseTableRow));
    return tbody;
}

export function createCourseTable(courses: Course[]): HTMLElement {
    const headers = ["Code", "Language", "Seen (words)", "Studied today"];
    const table = createTable(
        createTableHeader(headers),
        createCourseTableBody(courses),
    );
    return createScrollingTable(table);
}
