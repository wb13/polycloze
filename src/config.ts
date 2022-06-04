// Get location of server, but with a different port number.
function getLocation (port: string): string {
  const url = new URL(location.href)
  url.port = port
  return url.href
}

// Location of server
export const src = getLocation('3000')

let l1 = 'eng'
let l2 = 'spa'

export function getL1 (): string {
  return l1
}

export function getL2 (): string {
  return l2
}

export function setL1 (code: string) {
  l1 = code
}

export function setL2 (code: string) {
  l2 = code
}

export function currentCourse (): string {
  return `/${l1}/${l2}`
}
