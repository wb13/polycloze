import './index.css'
import { createApp } from './app'
import { fetchItems, submitReview } from './data'

createApp(fetchItems, [], submitReview)
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
