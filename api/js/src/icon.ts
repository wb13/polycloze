import "./icon.css";

// Creates SVG icon from ../public
export function createIcon(name: string): HTMLImageElement {
  const img = document.createElement("img");
  img.src = `/svg/ph@1.4.0/${name}.svg`;
  return img;
}

export function createLabeledIcon(
  name: string,
  label: string
): DocumentFragment {
  const icon = createIcon(name);
  const fragment = document.createDocumentFragment();
  fragment.append(icon);
  icon.alt = ""; // Because the icon is already described by the label.
  fragment.append(` ${label}`);
  return fragment;
}
