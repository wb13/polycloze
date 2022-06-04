import './select.css'
import { getL1, getL2, setL1, setL2 } from './data'

type Language = {
  code: string
  native: string
  english: string
}

function createLanguageOption (language: Language, selected: boolean = false): HTMLOptionElement {
  const option = document.createElement('option')
  option.value = language.code
  option.selected = selected
  option.textContent = language.native
  return option
}

// NOTE Doesn't update config.l1 or l2.
// You have to do it using the onChange callback function.
function createLanguageSelectInput (languages: Language[], selected: string = 'spa', onChange: () => void = () => {}): HTMLSelectElement {
  const select = document.createElement('select')
  select.append(...languages.map(language => createLanguageOption(language, language.code === selected)))
  select.addEventListener('change', onChange)
  return select
}

export function createLanguageSelectForm (languages: Language[]): HTMLFormElement {
  const form = document.createElement('form')

  const l1 = createLanguageSelectInput(languages, getL1(), () => {
    setL1(l1.value)
    location.reload()
  })
  const l2 = createLanguageSelectInput(languages.filter(l => l.code !== getL1()), getL2(), () => {
    setL2(l2.value)
    location.reload()
  })

  form.append('🌐 ', l1, ' > ', l2)
  return form
}
