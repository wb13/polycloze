/* Defines elements found in overview page (e.g. language stats). */

import { createButton } from "./button";
import { getL1, setL2 } from "./language";
import { Course, CourseStats } from "./schema";

function createButtonGroup(onClick: (event: Event) => void): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("button-group");
    p.appendChild(createButton("Start", onClick));
    return p;
}

function createSeenParagraph(stats: CourseStats): HTMLParagraphElement | null {
    const { seen, total } = stats;
    if (seen == null || seen < 0) {
        return null;
    }
    if (total == null || total <= 0) {
        return null;
    }

    const percentage = Math.floor(100 * seen / total);
    const text = `Seen: ${seen} word${seen === 1 ? "" : "s"} (${percentage}%)`;

    const p = document.createElement("p");
    p.textContent = text;
    return p;
}

function createTodayParagraph(stats: CourseStats): HTMLParagraphElement | null {
    const { learned, reviewed } = stats;
    const texts = [];
    if (learned != null && learned > 0) {
        texts.push(`${learned} new word${learned === 1 ? "" : "s"}`);
    }
    if (reviewed != null && reviewed > 0) {
        texts.push(`${reviewed} reviewed`);
    }

    if (texts.length === 0) {
        return null;
    }

    const p = document.createElement("p");
    p.textContent = `Today: ${texts.join(", ")}`;
    return p;
}

function createCourseStats(course: Course): HTMLDivElement | null {
    if (!course.stats) {
        return null;
    }

    const children = [];
    const seenParagraph = createSeenParagraph(course.stats);
    const todayParagraph = createTodayParagraph(course.stats);
    if (seenParagraph != null) {
        children.push(seenParagraph);
    }
    if (todayParagraph != null) {
        children.push(todayParagraph);
    }

    if (children.length === 0) {
        return null;
    }

    const div = document.createElement("div");
    div.append(...children);
    return div;
}

function createCourseOverview(course: Course, target: string): HTMLDivElement {
    const card = document.createElement("div");
    card.classList.add("card");
    card.innerHTML = `<h2>${course.l2.name}</h2>`;

    const row = document.createElement("div");
    row.classList.add("row");
    card.appendChild(row);

    const stats = createCourseStats(course);
    if (stats != null) {
        row.appendChild(stats);
    }

    const start = () => {
        setL2(course.l2);
        window.location.pathname = target;
    };
    row.appendChild(createButtonGroup(start));
    return card;
}

export function createOverview(courses: Course[], target = "/study"): HTMLDivElement {
    const l1 = getL1().code;
    const div = document.createElement("div");
    for (const course of courses) {
        if (course.l1.code === l1) {
            div.appendChild(createCourseOverview(course, target));
        }
    }
    return div;
}
