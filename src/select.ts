import './select.css'
import { getL1, getL2, setL1, setL2 } from './data'

export type LanguageStats = {
  seen?: number
  total?: number
  learned?: number
  reviewed?: number
};

export type Language = {
  code: string
  native: string
  english: string

  stats?: LanguageStats;
}

function createLanguageOption (language: Language, selected: boolean = false): HTMLOptionElement {
  const option = document.createElement('option')
  option.value = language.code
  option.selected = selected
  option.textContent = language.native
  return option
}

// NOTE Doesn't update l1 and l2 in localStorage.
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

/* Drop-down language select */

function createLanguageSelect (languages: Language[]): HTMLSelectElement {
  const select = document.createElement('select')
  select.id = 'language-select'

  for (const language of languages) {
    const option = createLanguageOption(language, language.code === getL1())
    select.appendChild(option)
  }
  return select
}

function createLanguageLabel (): HTMLLabelElement {
  const label = document.createElement('label')
  label.htmlFor = 'language-select'
  label.textContent = '🌐 '
  return label
}

export function createLanguageForm (languages: Language[]): HTMLFormElement {
  const form = document.createElement('form')
  form.append(createLanguageLabel(), createLanguageSelect(languages))
  return form
}
