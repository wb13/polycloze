import {
  fetchActivity,
  fetchCourses,
  fetchEstimatedLevel,
  fetchVocabularySize,
} from "./api";
import { createApp, createListenApp } from "./app";
import { ItemBuffer, RandomSentenceBuffer } from "./buffer";
import { setButtonLink } from "./button";
import { createScoreCounter } from "./counter";
import { createDiacriticButtonSettingsSection } from "./diacritic";
import { getL2 } from "./language";
import { createResponsiveMenu } from "./menu";
import { createOverviewPage } from "./overview";
import { ActivitySummary, Course, DataPoint } from "./schema";
import { createCourseSelectButton } from "./select";
import { createVoiceSettingsSection, getListenLevel, TTS } from "./tts";
import { createFileBrowser } from "./upload";
import { createVocabularyList } from "./vocab";

export class ClozeApp extends HTMLElement {
  promise: Promise<[HTMLDivElement, () => void]>;

  constructor() {
    super();

    const buffer = new ItemBuffer();
    this.promise = createApp(buffer);
  }

  async connectedCallback() {
    const [app, ready] = await this.promise;
    this.appendChild(app);
    ready();

    const l2 = getL2().name;
    document.title = `${l2} | polycloze`;
  }
}

export class ListenApp extends HTMLElement {
  promise: Promise<[HTMLDivElement, () => void]>;

  constructor() {
    super();

    const targetLevel = getListenLevel();
    const buffer = new RandomSentenceBuffer(targetLevel !== null ? targetLevel : 3);
    this.promise = createListenApp(buffer);
  }

  async connectedCallback() {
    //const estimatedLevel = await this.estimatedLevel;
    //const targetLevel = Math.max(estimatedLevel[estimatedLevel.length-1].value-1, 1);

    const [app, ready] = await this.promise;
    this.appendChild(app);
    ready();

    const l2 = getL2().name;
    document.title = `${l2} | polycloze`;
  }
}

export class CourseSelectButton extends HTMLElement {
  courses: Promise<Course[]>;

  constructor() {
    super();
    this.courses = fetchCourses();
  }

  async connectedCallback() {
    const courses = await this.courses;
    this.appendChild(createCourseSelectButton(courses));
  }
}

export class ResponsiveMenu extends HTMLElement {
  connectedCallback() {
    const signedIn = this.getAttribute("signed-in") != null;
    this.appendChild(createResponsiveMenu(signedIn));
  }
}

export class Overview extends HTMLElement {
  activity: Promise<ActivitySummary[]>;
  vocabularySize: Promise<DataPoint[]>;
  estimatedLevel: Promise<DataPoint[]>;

  constructor() {
    super();
    this.activity = fetchActivity();
    this.vocabularySize = fetchVocabularySize();
    this.estimatedLevel = fetchEstimatedLevel();
  }

  async connectedCallback() {
    const resolved = await Promise.all([
      this.activity,
      this.vocabularySize,
      this.estimatedLevel,
    ]);
    const [activity, vocabularySize, estimatedLevel] = resolved;
    const page = createOverviewPage(activity, vocabularySize, estimatedLevel);
    this.appendChild(page);
  }
}

export class ScoreCounter extends HTMLElement {
  activity: Promise<ActivitySummary[]>;

  constructor() {
    super();
    this.activity = fetchActivity();
  }

  async connectedCallback() {
    // TODO only fetch activity today
    const activity = await this.activity;
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
customElements.define("listen-app", ListenApp);
customElements.define("course-select-button", CourseSelectButton);
customElements.define("responsive-menu", ResponsiveMenu);
customElements.define("polycloze-overview", Overview);
customElements.define("score-counter", ScoreCounter);
customElements.define("button-link", ButtonLink, { extends: "button" });
customElements.define("vocabulary-list", VocabularyList);
customElements.define("course-settings", CourseSettings);
customElements.define("file-browser", FileBrowser);
