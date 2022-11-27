import "./select.css";
import { createButton } from "./button";
import { createIcon } from "./icon";
import { getL1, getL2, setL1, setL2 } from "./language";
import { createModal } from "./modal";
import { Course, Language } from "./schema";

function createLanguageSelect(selected: string, languages: Language[], onChange: (selected: string) => void): HTMLSelectElement {
    const select = document.createElement("select");
    select.addEventListener("change", () => onChange(select.value));

    for (const language of languages) {
        const option = document.createElement("option");
        option.value = language.code;
        option.textContent = language.name;
        if (selected === language.code) {
            option.selected = true;
        }
        select.appendChild(option);
    }
    return select;
}

// onChange: callback function when selection changes.
// It takes the language code of the selected language as input.
function createL1Select(l1Code: string, languages: Language[], onChange: (selected: string) => void): HTMLParagraphElement {
    const select = createLanguageSelect(l1Code, languages, onChange);
    select.id = "l1-select";

    const label = document.createElement("label");
    label.htmlFor = select.id;
    label.textContent = "I speak ";

    const p = document.createElement("p");
    p.append(label, document.createElement("br"), select);
    return p;
}

function createL2Select(l2Code: string, languages: Language[], onChange: (selected: string) => void): HTMLParagraphElement {
    const select = createLanguageSelect(l2Code, languages, onChange);
    select.id = "l2-select";

    const label = document.createElement("label");
    label.htmlFor = select.id;
    label.textContent = "I want to learn ";

    const p = document.createElement("p");
    p.append(label, document.createElement("br"), select);
    return p;
}

function createSaveButtonGroup(onClick: (event: Event) => void): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("button-group");

    const content = document.createElement("span");
    content.append(createIcon("floppy-disk"), " Save");

    const button = createButton(content, onClick);
    p.appendChild(button);
    return p;
}

function compareLanguages(a: Language, b: Language): number {
    if (a.code < b.code) {
        return -1;
    } else if (a.code === b.code) {
        return 0;
    } else {
        return 1;
    }
}

function sourceLanguages(courses: Course[]): Language[] {
    const visited = new Set();
    const languages = [];
    for (const { l1 } of courses) {
        if (visited.has(l1.code)) {
            continue;
        }
        visited.add(l1.code);
        languages.push(l1);
    }
    return languages.sort(compareLanguages);
}

function targetLanguages(l1Code: string, courses: Course[]): Language[] {
    const visited = new Set();
    const languages = [];
    for (const course of courses) {
        if (l1Code !== course.l1.code || visited.has(course.l2.code)) {
            continue;
        }
        visited.add(course.l2.code);
        languages.push(course.l2);
    }
    return languages.sort(compareLanguages);
}

function languageCodes(courses: Course[]): Map<string, Language> {
    const codes = new Map();
    for (const course of courses) {
        codes.set(course.l1.code, course.l1);
        codes.set(course.l2.code, course.l2);
    }
    return codes;
}

// Creates a course select menu/modal
function createCourseMenu(courses: Course[]): [HTMLDivElement, () => void] {
    const codes = languageCodes(courses);
    let l1 = getL1();
    let l2 = getL2();

    const selectL1 = createL1Select(l1.code, sourceLanguages(courses), updateL1);
    let selectL2 = createL2Select(l2.code, targetLanguages(l1.code, courses), updateL2);

    const div = document.createElement("div");
    div.classList.add("course-menu");
    div.append(selectL1, selectL2, createSaveButtonGroup(save));
    return createModal(div);

    function updateL1(code: string) {
        if (code === l2.code) {
            l2 = l1;
        }
        l1 = codes.get(code) as Language;
        updateChoices();
    }

    function updateL2(code: string) {
        l2 = codes.get(code) as Language;
        updateChoices();
    }

    function updateChoices() {
        const languages = targetLanguages(l1.code, courses);
        const replacement = createL2Select(l2.code, languages, updateL2);
        selectL2.replaceWith(replacement);
        selectL2 = replacement;
    }

    function save() {
        setL1(l1);
        setL2(l2);
        location.reload();
    }
}

export function createCourseSelectButton(courses: Course[]): HTMLButtonElement {
    const [menu, show] = createCourseMenu(courses);
    document.body.appendChild(menu);

    const l1 = getL1();
    const l2 = getL2();
    const content = document.createElement("span");
    content.append(createIcon("translate"), ` ${l1.code}-${l2.code}`);

    const button = createButton(content, show);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}
