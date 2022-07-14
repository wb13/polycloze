// Get location of server, but with a different port number.
function getLocation(port: string): string {
    const url = new URL(location.href);
    url.port = port;
    return url.href;
}

// Location of server
export const src = getLocation("3000");
