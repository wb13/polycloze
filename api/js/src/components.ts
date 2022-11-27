import { fetchActivityHistory, fetchCourses } from "./api";
import { createApp } from "./app";
import { ItemBuffer } from "./buffer";
import { setButtonLink } from "./button";
import { createScoreCounter } from "./counter";
import { getL2 } from "./language";
import { createResponsiveMenu } from "./menu";
import { createOverviewPage } from "./overview";
import { createCourseSelectButton } from "./select";
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

export class CourseSelectButton extends HTMLElement {
    async connectedCallback() {
        const courses = await fetchCourses();
        this.appendChild(createCourseSelectButton(courses));
    }
}

export class ResponsiveMenu extends HTMLElement {
    async connectedCallback() {
        const signedIn = this.getAttribute("signed-in") != null;
        this.appendChild(createResponsiveMenu(signedIn));
    }
}

export class Overview extends HTMLElement {
    async connectedCallback() {
        const activityHistory = await fetchActivityHistory();
        this.appendChild(createOverviewPage(activityHistory));
    }
}

export class ScoreCounter extends HTMLElement {
    async connectedCallback() {
        const activityHistory = await fetchActivityHistory();
        const today = activityHistory.activities[0];
        const { crammed, learned, strengthened } = today;
        const score = crammed + learned + strengthened;
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
customElements.define("course-select-button", CourseSelectButton);
customElements.define("responsive-menu", ResponsiveMenu);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
customElements.define("button-link", ButtonLink, { extends: "button" });
customElements.define("vocabulary-list", VocabularyList);
