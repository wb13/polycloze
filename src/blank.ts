import './blank.css'
import { distance } from 'fastest-levenshtein'
import { getFont, getWidth, substituteDigraphs } from './text'

type Status = 'correct' | 'incorrect' | 'almost';

function changeStatus (input: HTMLInputElement, status: Status) {
  input.classList.add(status)
}

// enable: Enable submit button
// Also returns a resize function, which should be called when the element is
// connected to the DOM.
export function createBlank (answer: string, done: (answer: string, correct: true) => void, enable: () => void): [HTMLInputElement, () => void] {
  let correct = true
  const input = document.createElement('input')
  input.autocapitalize = 'none'
  input.classList.add('blank')

  input.addEventListener('input', () => {
    if (input.value !== '') {
      enable()
    }
    input.value = substituteDigraphs(input.value)
  })
  input.addEventListener('change', () => {
    switch (distance(input.value, answer)) {
      case 0:
        changeStatus(input, 'correct')
        return done(answer, correct)

      case 1:
      case 2:
        changeStatus(input, 'almost')
        break

      default:
        correct = false
        input.placeholder = answer
        input.value = ''
        changeStatus(input, 'incorrect')
        break
    }
  })

  const resize = () => {
    const width = getWidth(getFont(input), answer)
    input.style.setProperty('width', width)
  }
  return [input, resize]
}
