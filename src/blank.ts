import './blank.css'

import { distance } from 'fastest-levenshtein'

type Status = 'correct' | 'incorrect' | 'almost';

function changeStatus (input: HTMLInputElement, status: Status) {
  input.classList.add(status)
}

// enable: Enable submit button
export function createBlank (answer: string, done: (correct: true) => void, enable: () => void): HTMLInputElement {
  let correct = true
  const input = document.createElement('input')
  input.classList.add('blank')
  input.addEventListener('input', () => {
    if (input.value !== '') {
      enable()
    }
  })
  input.addEventListener('change', () => {
    switch (distance(input.value, answer)) {
      case 0:
        changeStatus(input, 'correct')
        return done(correct)

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
  return input
}
