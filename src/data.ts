import { Item } from './item'

export async function fetchItems (src: string): Promise<Item[]> {
  const request = new Request(src, { mode: 'cors' })
  const response = await fetch(request)
  const json = await response.json()
  return json.items
}

type Review = {
  word: string
  correct: boolean
}

export async function submitReviews (src: string, reviews: Review[]): Promise<boolean> {
  const options = {
    body: JSON.stringify({ reviews }),
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
