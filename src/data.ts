import { src } from './config'
import { Item } from './item'
import { Language } from './select'

async function fetchJson (url: string, options: any): Promise<any> {
  const request = new Request(url, options)
  const response = await fetch(request)
  return await response.json()
}

export async function supportedLanguages (): Promise<Language[]> {
  const url = new URL('/options', src)
  const options = { mode: 'cors' }
  const json = await fetchJson(url, options)
  return json.languages
}

export async function fetchItems (): Promise<Item[]> {
  const json = await fetchJson(src, { mode: 'cors' })
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
  return await fetchJson(src, options).success
}
