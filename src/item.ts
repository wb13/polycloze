import './item.css'
import { createButton } from './button'
import { Sentence, createSentence } from './sentence'

type Item = {
  sentence: Sentence
  translation: string
}

function createItemHeader (): HTMLDivElement {
  // TODO
  return document.createElement('div')
}

function createTranslation (translation: string): HTMLDivElement {
  const div = document.createElement('div')
  div.classList.add('translation')
  div.textContent = translation
  return div
}

function createItemBody (item: Item, next: (ok: boolean) => void, enable: () => void, post: (word: string, correct: boolean) => void): [HTMLDivElement, () => void, () => void] {
  const div = document.createElement('div')
  const [sentence, check, resize] = createSentence(item.sentence, next, enable, post)
  div.append(
    sentence,
    createTranslation(item.translation)
  )
  return [div, check, resize]
}

function createItemFooter (submitBtn: HTMLButtonElement): HTMLDivElement {
  const div = document.createElement('div')
  div.classList.add('button-group')
  div.appendChild(submitBtn)
  return div
}

function createSubmitButton (onClick: (event: Event) => void): [HTMLButtonElement, () => void] {
  const button = createButton('Submit', onClick)
  button.disabled = true

  const enable = () => {
    button.disabled = false
  }
  return [button, enable]
}

export function createItem (item: Item, next: (ok: boolean) => void, post: (word: string, correct: boolean) => void): [HTMLDivElement, () => void] {
  const [submitBtn, enable] = createSubmitButton()

  const header = createItemHeader()
  const [body, check, resize] = createItemBody(item, next, enable, post)
  const footer = createItemFooter(submitBtn)

  submitBtn.addEventListener('click', check)

  const div = document.createElement('div')
  div.classList.add('item')
  div.append(header, body, footer)
  return [div, resize]
}
