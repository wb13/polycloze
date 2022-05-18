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
  div.textContent = translation
  return div
}

function createItemBody (item: Item, next: (ok: boolean) => void, enable: () => void): [HTMLDivElement, () => void] {
  const div = document.createElement('div')
  const [sentence, check] = createSentence(item.sentence, next, enable)
  div.append(
    sentence,
    createTranslation(item.translation)
  )
  return [div, check]
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

export function createItem (item: Item, next: (ok: boolean) => void): HTMLDivElement {
  const [submitBtn, enable] = createSubmitButton()

  const header = createItemHeader()
  const [body, check] = createItemBody(item, next, enable)
  const footer = createItemFooter(submitBtn)

  submitBtn.addEventListener('click', check)

  const div = document.createElement('div')
  div.classList.add('item')
  div.append(header, body, footer)
  return div
}

export function createFromItems (items: Item[]): [HTMLDivElement, () => void] {
  if (items.length === 0) {
    throw new Error('unhandled case')
  }

  const div = document.createElement('div')
  const item = items.pop()
  const next = () => {
    const [replacement, ready] = createFromItems(items)
    div.replaceWith(replacement)
    ready()
  }
  div.appendChild(createItem(item, next))

  const ready = () => {
    div.querySelector('.blank')?.focus()
  }
  return [div, ready]
}
