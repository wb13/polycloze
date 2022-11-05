/* Course table. */

import "./course.css";
import { createButton } from "./button";
import { createIcon } from "./icon";
import { setL2 } from "./language";
import { Course } from "./schema";
import { createTable, createTableData, createTableHeader } from "./table";

function createSpan(children: Array<string | Element>): HTMLSpanElement {
    const span = document.createElement("span");
    span.append(...children);
    return span;
}

function createSeenCount(course: Course): HTMLElement {
    const seen = course.stats?.seen || 0;
    if (seen <= 0) {
        return document.createElement("div");
    }

    const content = createSpan([createIcon("info"), ` ${seen} words`]);
    const button = createButton(content, () => {
        setL2(course.l2);
        window.location.pathname = "/vocab";
    });
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    button.style.margin = "0";
    return button;
}

function createDiv(child: string | Element): HTMLDivElement {
    const div = document.createElement("div");
    div.append(child);
    return div;
}

function createCourseCodeButton(course: Course): HTMLButtonElement {
    const content = createSpan([createIcon("play-circle"), ` ${course.l2.code}`]);
    const button = createButton(content, () => {
        setL2(course.l2);
        window.location.pathname = "/study";
    });
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    button.style.margin = "0";
    return button;
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
        createTableData(createCourseCodeButton(course)),
        createTableData(course.l2.name),
        createTableData(createSeenCount(course)),
        createStudiedTodayCell(course),
    );
    return tr;
}

function createCourseTableBody(courses: Course[]): HTMLTableSectionElement {
    const tbody = document.createElement("tbody");
    tbody.append(...courses.map(createCourseTableRow));
    return tbody;
}

export function createCourseTable(courses: Course[]): HTMLTableElement {
    const headers = ["Code", "Language", "Seen (words)", "Studied today"];
    return createTable(
        createTableHeader(headers),
        createCourseTableBody(courses),
    );
}
