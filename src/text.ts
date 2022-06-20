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

  ['~n', 'ñ'],
  ['~N', 'Ñ'],

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

  ['ss', 'ß']
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
