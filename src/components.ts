import { createApp } from './app'
import { ItemBuffer } from './buffer'
import { supportedLanguages } from './data'
import { createLanguageForm, createLanguageSelectForm } from './select'

export class ClozeApp extends HTMLElement {
  async connectedCallback () {
    const languages = await supportedLanguages()
    const form = createLanguageSelectForm(languages)
    this.appendChild(form)

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

customElements.define('cloze-app', ClozeApp)
customElements.define('language-select', LanguageSelect)
