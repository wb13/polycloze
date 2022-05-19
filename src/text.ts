const canvas = document.createElement('canvas')

export function getFont (ele: Element): string {
  const style = getComputedStyle(ele)
  const weight = style.getPropertyValue('font-weight') || 'normal'
  const size = style.getPropertyValue('font-size') || '12px'
  const family = style.getPropertyValue('font-family') || 'sans-serif'
  return `${weight} ${size} ${family}`
}

export function getWidth (font: string, text: string): string {
  const context = canvas.getContext('2d')
  context.font = font
  const metrics = context.measureText(text)
  return `${metrics.width}px`
}

const digraphs = new Map([
  ["'a", 'á'],
  ["'e", 'é'],
  ["'i", 'í'],
  ["'o", 'ó'],
  ["'u", 'ú'],
  ["'A", 'Á'],
  ["'E", 'É'],
  ["'I", 'Í'],
  ["'O", 'Ó'],
  ["'U", 'Ú'],

  [':a', 'ä'],
  [':e', 'ë'],
  [':i', 'ï'],
  [':o', 'ö'],
  [':u', 'ü'],
  [':A', 'Ä'],
  [':E', 'Ë'],
  [':I', 'Ï'],
  [':O', 'Ö'],
  [':U', 'Ü'],

  ['~n', 'ñ'],
  ['~N', 'Ñ']
])

function reverseString (text: string): string {
  return text.split('').reverse().join('')
}

function substituteDigraph (digraph: string): string {
  if (!digraph.startsWith('\\')) {
    throw new Error("digraph should start with '\\'")
  }

  const key = digraph.slice(1)
  const result = digraphs.get(key)
  if (result != null) {
    return result
  }
  return digraphs.get(reverseString(key)) || digraph
}

export function substituteDigraphs (text: string): string {
  let start = 0
  while (start < text.length) {
    if (text[start] === '\\') {
      const digraph = text.slice(start, start + 3)
      text = text.slice(0, start) + substituteDigraph(digraph) + text.slice(start + 3)
    }
    ++start
  }
  return text
}
