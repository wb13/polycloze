import './index.css'
import { createApp } from './app'

createApp([])
  .then(([app, ready]) => {
    document.body.appendChild(app)
    ready()
  })
