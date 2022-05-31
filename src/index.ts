import './index.css'
import { createApp } from './app'
import { fetchItems, submitReviews } from './data'

const post = (word: string, correct: boolean) => submitReviews([{ word, correct }])

createApp(fetchItems, [], post)
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
