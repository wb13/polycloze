import './sentence.css'
import { createBlank } from './blank'
import { submitReview } from './data'

export type Sentence = {
  id: number
  parts: string[]
}

function createPart (part: string): HTMLSpanElement {
  const span = document.createElement('span')
  span.textContent = part
  return span
}

function isBeginning (part: string): boolean {
  switch (part) {
    case '':
    case '¿':
    case '¡':
      return true

    default:
      return false
  }
}

// TODO document params
// Note: takes two callback functions.
// - next: ?
// - enable: Enables submit button (used by createBlank).
//
// In addition to a div element, also returns two functions to be called by the
// caller.
// - check: ?
// - resize: ?
export function createSentence (sentence: Sentence, next: (ok: boolean) => void, enable: () => void): [HTMLDivElement, () => void, () => void] {
  let ok = true
  let remaining = Math.floor(sentence.parts.length / 2)
  const check = () => {
    if (remaining <= 0) {
      next(ok)
    }
  }
  const done = (answer: string, correct: boolean) => {
    submitReview(answer, correct)
    --remaining
    if (!correct) {
      ok = false
    }
    check()
  }

  const resizeFns: Array<() => void> = []

  const div = document.createElement('div')
  div.classList.add('sentence')
  for (const [i, part] of sentence.parts.entries()) {
    if (i % 2 === 0) {
      div.appendChild(createPart(part))
    } else {
      const autocapitalize = (i === 1) && isBeginning(sentence.parts[0])
      const [blank, resize] = createBlank(part, autocapitalize, done, enable)
      div.appendChild(blank)
      resizeFns.push(resize)
    }
  }

  fixPunctuationWrap(div)

  const resizeAll = () => {
    for (const fn of resizeFns) {
      fn()
    }
  }
  return [div, check, resizeAll]
}

// Prevents punctuation symbols from starting a new line.
// Assumes all child nodes are elements.
function fixPunctuationWrap (div: HTMLDivElement) {
  const inputs = div.querySelectorAll('.blank')
  for (let i = 0; i < inputs.length; i++) {
    const input = inputs[i]
    const span = input.nextElementSibling

    if (span == null) {
      continue
    }

    // NOTE Does not split by other whitespace characters
    const words = span!.textContent?.split(' ') || []
    if (words.length > 0 && words[0] !== '') {
      const wrapper = document.createElement('span')
      wrapper.style.whiteSpace = 'nowrap'
      input.replaceWith(wrapper)
      wrapper.appendChild(input)

      const after = document.createElement('span')
      after.textContent = words[0]
      wrapper.appendChild(after)

      words.shift()

      const tail = words.join(' ')
      span.textContent = tail.length > 0 ? ' ' + tail : ''
    }
  }
}
