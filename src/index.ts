import { fetchItems } from './data'
import { createFromItems } from './item'

const src = 'http://localhost:3000'
fetchItems(src)
  .then(items => {
    const [div, ready] = createFromItems(items)
    document.body.appendChild(div)
    ready()
  })
