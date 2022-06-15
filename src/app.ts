import './app.css'
import { bufferedFetchItems } from './data'
import { Item, createItem } from './item'

export async function createApp (items: Item[]): Promise<[HTMLDivElement, () => void]> {
  if (items.length === 0) {
    return createApp(await bufferedFetchItems())
  }

  const div = document.createElement('div')
  const item = items.pop()
  const next = () => {
    createApp(items).then(([replacement, ready]) => {
      div.replaceWith(replacement)
      ready()
    })
  }

  const [child, resize] = createItem(item, next)
  div.appendChild(child)

  const ready = () => {
    div.querySelector('.blank')?.focus()
    resize()
  }
  return [div, ready]
}
