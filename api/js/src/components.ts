import {
  fetchActivity,
  fetchCourses,
  fetchEstimatedLevel,
  fetchVocabularySize,
} from "./api";
import { createApp } from "./app";
import { ItemBuffer } from "./buffer";
import { setButtonLink } from "./button";
import { createScoreCounter } from "./counter";
import { createDiacriticButtonSettingsSection } from "./diacritic";
import { getL2 } from "./language";
import { createResponsiveMenu } from "./menu";
import { createOverviewPage } from "./overview";
import { createCourseSelectButton } from "./select";
import { createVoiceSettingsSection, TTS } from "./tts";
import { createFileBrowser } from "./upload";
import { createVocabularyList } from "./vocab";

export class ClozeApp extends HTMLElement {
  async connectedCallback() {
    const l2 = getL2().name;
    const buffer = new ItemBuffer();
    const [app, ready] = await createApp(buffer);
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
    const resolved = await Promise.all([
      fetchActivity(),
      fetchVocabularySize(),
      fetchEstimatedLevel(),
    ]);
    const [activity, vocabularySize, estimatedLevel] = resolved;
    const page = createOverviewPage(activity, vocabularySize, estimatedLevel);
    this.appendChild(page);
  }
}

export class ScoreCounter extends HTMLElement {
  async connectedCallback() {
    // TODO only fetch activity today
    const activity = await fetchActivity();
    const today = activity[activity.length - 1];
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
  tts: TTS;
  init: Promise<void>;
  // Await `this.init` to make sure `tts` is initialized.

  constructor() {
    super();
    this.tts = new TTS();
    this.init = this.tts.init();
  }

  async connectedCallback() {
    await this.init;
    this.appendChild(await createVocabularyList(this.tts));
  }
}

export class CourseSettings extends HTMLElement {
  tts: TTS;
  init: Promise<void>;
  // Await `this.init` to make suer `tts` is initialized.

  constructor() {
    super();
    this.tts = new TTS();
    this.init = this.tts.init();
  }

  async connectedCallback() {
    await this.init;
    this.append(
      createDiacriticButtonSettingsSection(),
      document.createElement("br"),
      createVoiceSettingsSection(this.tts)
    );
  }
}

export class FileBrowser extends HTMLElement {
  connectedCallback() {
    const name = this.getAttribute("name") || "csv-upload";
    this.appendChild(createFileBrowser(name));
  }
}

customElements.define("cloze-app", ClozeApp);
customElements.define("course-select-button", CourseSelectButton);
customElements.define("responsive-menu", ResponsiveMenu);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
customElements.define("button-link", ButtonLink, { extends: "button" });
customElements.define("vocabulary-list", VocabularyList);
customElements.define("course-settings", CourseSettings);
customElements.define("file-browser", FileBrowser);
