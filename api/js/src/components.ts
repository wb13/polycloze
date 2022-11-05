import { createApp } from "./app";
import { ItemBuffer } from "./buffer";
import { setButtonLink } from "./button";
import { createScoreCounter } from "./counter";
import { createCourseTable } from "./course";
import { availableCourses } from "./data";
import { getL1, getL2 } from "./language";
import { createLanguageSelectButton } from "./select";
import { createVocabularyList } from "./vocab";

export class ClozeApp extends HTMLElement {
    async connectedCallback() {
        const l2 = getL2().name;
        const [app, ready] = await createApp(new ItemBuffer());
        this.appendChild(app);
        ready();
        document.title = `${await l2} | polycloze`;
    }
}

export class LanguageSelectButton extends HTMLElement {
    async connectedCallback() {
        const courses = await availableCourses();
        const languages = courses.map(c => c.l1);
        this.appendChild(createLanguageSelectButton(languages));
    }
}

export class Overview extends HTMLElement {
    async connectedCallback() {
        const courses = await availableCourses({
            l1: getL1().code,
            stats: true,
        });
        this.innerHTML = "<h1>Pick a language.</h1>";
        this.appendChild(createCourseTable(courses));
    }
}

export class ScoreCounter extends HTMLElement {
    async connectedCallback() {
        const courses = await availableCourses({
            l1: getL1().code,
            l2: getL2().code,
            stats: true,
        });

        const l1 = getL1().code;
        const l2 = getL2().code;

        const course = courses.find(c => c.l1.code === l1 && c.l2.code === l2);
        const score = course?.stats?.correct || 0;
        this.appendChild(createScoreCounter(score));
    }
}

export class ButtonLink extends HTMLButtonElement {
    connectedCallback() {
        const href = this.getAttribute("href") || "/";
        const method = (this.getAttribute("method") || "GET").toUpperCase();
        setButtonLink(this, href, method);
    }
}

export class VocabularyList extends HTMLElement {
    async connectedCallback() {
        this.appendChild(await createVocabularyList());
    }
}

customElements.define("cloze-app", ClozeApp);
customElements.define("language-select-button", LanguageSelectButton);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
customElements.define("button-link", ButtonLink, { extends: "button" });
customElements.define("vocabulary-list", VocabularyList);
