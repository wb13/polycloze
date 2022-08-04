import { createApp } from "./app";
import { ItemBuffer } from "./buffer";
import { createScoreCounter } from "./counter";
import { getL1, getL2, availableCourses } from "./data";
import { createOverview } from "./overview";
import { createLanguageForm } from "./select";

export class ClozeApp extends HTMLElement {
    async connectedCallback() {
        const l2 = this.getL2Name();

        const [app, ready] = await createApp(new ItemBuffer());
        this.appendChild(app);
        ready();
        document.title = `polycloze | ${await l2}`;
    }

    async getL2Name(): Promise<string> {
        const code = getL2();
        const courses = await availableCourses();
        const course = courses.find(c => c.l2.code === code);
        return course ? course.l2.name : code;
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
        const target = this.getAttribute("target") || "/study";
        const courses = await availableCourses();
        this.innerHTML = "<h1>Pick a language.</h1>";
        this.appendChild(createOverview(courses, target));
    }
}

export class ScoreCounter extends HTMLElement {
    async connectedCallback() {
        const courses = await availableCourses();
        const course = courses.find(c => c.l1.code === getL1() && c.l2.code === getL2());
        const score = course?.stats?.correct || 0;
        this.appendChild(createScoreCounter(score));
    }
}

customElements.define("cloze-app", ClozeApp);
customElements.define("language-select", LanguageSelect);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
