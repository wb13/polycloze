import './index.css'
import { createApp } from './app'
import { supportedLanguages } from './data'
import { createLanguageSelectForm } from './select'

window.onload = () => {
  supportedLanguages().then(languages => {
    const form = createLanguageSelectForm(languages)
    document.body.appendChild(form)
  })
}

createApp([])
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
