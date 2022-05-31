import { src } from './config'
import { Item } from './item'

export async function fetchItems (): Promise<Item[]> {
  const request = new Request(src, { mode: 'cors' })
  const response = await fetch(request)
  const json = await response.json()
  return json.items
}

// Returns response status (success or not).
export async function submitReview (word: string, correct: boolean): Promise<boolean> {
  const options = {
    body: JSON.stringify({
      reviews: [
        { word, correct }
      ]
    }),
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    method: 'POST',
    mode: 'cors'
  }
  const request = new Request(src, options)
  const response = await fetch(request)
  const json = await response.json()
  return json.success
}
