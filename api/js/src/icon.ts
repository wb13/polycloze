import "./icon.css";

// Creates SVG icon from ../public
export function createIcon(name: string): HTMLImageElement {
    const img = document.createElement("img");
    img.src = `/public/svg/${name}.svg`;
    return img;
}
