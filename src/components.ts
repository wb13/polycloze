import { createApp } from './app'
import { ItemBuffer } from './buffer'
import { supportedLanguages } from './data'
import { createLanguageSelectForm } from './select'

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

customElements.define('cloze-app', ClozeApp)
