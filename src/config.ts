// Compile-time config

// Get location of server, but with a different port number.
function getLocation (port: string): string {
  const { protocol, hostname, pathname, search, hash } = location
  return `${protocol}//${hostname}:${port}${pathname}${search}${hash}`
}

export const src = getLocation('3000')
