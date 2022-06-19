import { createApp } from './app'
import { ItemBuffer } from './buffer'
import { createScoreCounter } from './counter'
import { getL2, supportedLanguages } from './data'
import { createOverview } from './overview'
import { createLanguageForm } from './select'

export class ClozeApp extends HTMLElement {
  async connectedCallback () {
    const [app, ready] = await createApp(new ItemBuffer())
    this.appendChild(app)
    ready()
  }
}

export class LanguageSelect extends HTMLElement {
  async connectedCallback () {
    const languages = await supportedLanguages()
    this.appendChild(createLanguageForm(languages))
  }
}

export class Overview extends HTMLElement {
  async connectedCallback () {
    const languages = await supportedLanguages()
    this.innerHTML = '<h1>Pick a language.</h1>'
    this.appendChild(createOverview(languages))
  }
}

export class ScoreCounter extends HTMLElement {
  async connectedCallback () {
    const languages = await supportedLanguages()
    const { stats } = languages.find(l => l.code === getL2())!

    this.appendChild(createScoreCounter(stats?.correct || 0, stats?.incorrect || 0))
  }
}

customElements.define('cloze-app', ClozeApp)
customElements.define('language-select', LanguageSelect)
customElements.define('polycloze-overview', Overview)
customElements.define('score-counter', ScoreCounter)
