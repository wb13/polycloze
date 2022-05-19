import './sentence.css'
import { createBlank } from './blank'

export type Sentence = {
  id: number
  parts: string[]
}

function createPart (part: string): HTMLSpanElement {
  const span = document.createElement('span')
  span.textContent = part
  return span
}

export function createSentence (sentence: Sentence, next: (ok: boolean) => void, enable: () => void): [HTMLDivElement, () => void, () => void] {
  let ok = true
  let remaining = Math.floor(sentence.parts.length / 2)
  const check = () => {
    if (remaining <= 0) {
      next(ok)
    }
  }
  const done = (correct: boolean) => {
    --remaining
    if (!correct) {
      ok = false
    }
    check()
  }

  const resizeFns = []

  const div = document.createElement('div')
  div.classList.add('sentence')
  for (const [i, part] of sentence.parts.entries()) {
    if (i % 2 === 0) {
      div.appendChild(createPart(part))
    } else {
      const [blank, resize] = createBlank(part, done, enable)
      div.appendChild(blank)
      resizeFns.push(resize)
    }
  }

  const resizeAll = () => {
    for (const fn of resizeFns) {
      fn()
    }
  }
  return [div, check, resizeAll]
}
