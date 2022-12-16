import "./icon.css";

// Creates SVG icon from ../public
export function createIcon(name: string): HTMLImageElement {
  const img = document.createElement("img");
  img.src = `/public/svg/${name}.svg?t=20221114`;
  // Update t value to bust cache (e.g. new version of icon library).
  return img;
}

export function createLabeledIcon(
  name: string,
  label: string
): DocumentFragment {
  const fragment = document.createDocumentFragment();
  fragment.append(createIcon(name));
  if (label) {
    fragment.append(` ${label}`);
  }
  return fragment;
}
