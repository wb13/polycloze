import './index.css'
import { createApp } from './app'
import { fetchItems } from './data'

const src = 'http://localhost:3000'
const refresh = () => fetchItems(src)

createApp(refresh, [])
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
