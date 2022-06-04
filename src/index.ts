import './index.css'
import { createApp } from './app'
import { supportedLanguages } from './data'
import { createLanguageSelectForm } from './select'

window.onload = async () => {
  const languages = await supportedLanguages()
  const form = createLanguageSelectForm(languages)
  document.body.appendChild(form)

  const [app, ready] = await createApp([])
  document.body.appendChild(app)
  ready()
}
