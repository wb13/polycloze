import './app.css'
import { ItemBuffer } from './buffer'
import { createItem } from './item'

export async function createApp (buffer: ItemBuffer): Promise<[HTMLDivElement, () => void]> {
  const div = document.createElement('div')
  const [item, afterCheck] = await buffer.take()
  const next = () => {
    createApp(buffer).then(([replacement, ready]) => {
      div.replaceWith(replacement)
      ready()
    })

    // NOTE next gets called in createSentence
    afterCheck()
  }

  const [child, resize] = createItem(item, next)
  div.appendChild(child)

  const ready = () => {
    div.querySelector('.blank')?.focus()
    resize()
  }
  return [div, ready]
}
