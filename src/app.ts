import './app.css'
import { Item, createItem } from './item'

export async function createApp (items: Item[], refresh: Promise<Item[]>, post: (word: string, correct: boolean) => void): Promise<[HTMLDivElement, () => void]> {
  if (items.length === 0) {
    return createApp(await refresh(), refresh, post)
  }

  const div = document.createElement('div')
  const item = items.pop()
  const next = () => {
    createApp(items, refresh, post).then(([replacement, ready]) => {
      div.replaceWith(replacement)
      ready()
    })
  }

  const [child, resize] = createItem(item, next, post)
  div.appendChild(child)

  const ready = () => {
    div.querySelector('.blank')?.focus()
    resize()
  }
  return [div, ready]
}
