import "./icon.css";

// Creates SVG icon from ../public
export function createIcon(name: string): HTMLImageElement {
    const img = document.createElement("img");
    img.src = `/public/svg/${name}.svg?t=20221114`;
    // Update t value to bust cache.
    return img;
}
