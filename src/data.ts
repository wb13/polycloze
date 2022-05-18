import { Item } from './item'

export async function fetchItems (src: string): Promise<Item[]> {
  const request = new Request(src, { mode: 'cors' })
  const response = await fetch(request)
  const json = await response.json()
  return json.items
}
