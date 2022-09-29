import { createApp } from "./app";
import { ItemBuffer } from "./buffer";
import { createScoreCounter } from "./counter";
import { availableCourses } from "./data";
import { getL1, getL2 } from "./language";
import { createLanguageForm } from "./select";
import { createCourseTable } from "./table";

export class ClozeApp extends HTMLElement {
    async connectedCallback() {
        const l2 = getL2().name;
        const [app, ready] = await createApp(new ItemBuffer());
        this.appendChild(app);
        ready();
        document.title = `polycloze | ${await l2}`;
    }
}

export class LanguageSelect extends HTMLElement {
    async connectedCallback() {
        const courses = await availableCourses();
        const languages = courses.map(c => c.l1);
        this.appendChild(createLanguageForm(languages));
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

customElements.define("cloze-app", ClozeApp);
customElements.define("language-select", LanguageSelect);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
