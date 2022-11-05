/* Course table. */

import "./course.css";
import { setL2 } from "./language";
import { Course } from "./schema";
import { createTable, createTableData, createTableHeader } from "./table";

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
        createTableData(course.l2.code),
        createTableData(course.l2.name),
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
    const headers = ["Code", "Language", "Seen (words)", "Studied today"];
    return createTable(
        createTableHeader(headers),
        createCourseTableBody(courses),
    );
}
