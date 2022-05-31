import './index.css'
import { createApp } from './app'
import { fetchItems, submitReviews } from './data'

const src = 'http://localhost:3000'
const refresh = () => fetchItems(src)

const post = (word: string, correct: boolean) => submitReviews(src, [
  { word, correct }
])

createApp(refresh, [], post)
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
