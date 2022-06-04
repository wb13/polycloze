// Contains functions for getting data from the server and from localStorage.

import { src } from './config'
import { Item } from './item'
import { Language } from './select'

// Local storage stuff

export function getL1 (): string {
  return localStorage.getItem('l1') || 'eng'
}

export function getL2 (): string {
  return localStorage.getItem('l2') || 'spa'
}

export function setL1 (code: string) {
  localStorage.setItem('l1', code)
}

export function setL2 (code: string) {
  localStorage.setItem('l2', code)
}

function currentCourse (): string {
  return `/${getL1()}/${getL2()}`
}

// Server stuff

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
  const url = new URL(currentCourse(), src)
  const options = { mode: 'cors' }
  const json = await fetchJson(url, options)
  return json.items
}

// Returns response status (success or not).
export async function submitReview (word: string, correct: boolean): Promise<boolean> {
  const url = new URL(currentCourse(), src)
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
  return await fetchJson(url, options).success
}
