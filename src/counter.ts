/* Score counter. */

// Returns a score counter in a div.
export function createScoreCounter (correct: number, incorrect: number): HTMLDivElement {
  const div = document.createElement('div')
  div.classList.add('medium')
  div.innerHTML = `<span class="correct">${correct}</span>/${correct + incorrect}`

  const listener = (event: Event) => {
    window.removeEventListener('polycloze-review', listener)
    if ((event as CustomEvent).detail.correct) {
      div.replaceWith(createScoreCounter(correct + 1, incorrect))
    } else {
      div.replaceWith(createScoreCounter(correct, incorrect + 1))
    }
  }

  window.addEventListener('polycloze-review', listener)
  return div
}
