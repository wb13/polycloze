import './index.css'
import { createApp } from './app'
import { ItemBuffer } from './buffer'
import { supportedLanguages } from './data'
import { createLanguageSelectForm } from './select'

window.onload = async () => {
  const languages = await supportedLanguages()
  const form = createLanguageSelectForm(languages)
  document.body.appendChild(form)

  const [app, ready] = await createApp(new ItemBuffer())
  document.body.appendChild(app)
  ready()
}
