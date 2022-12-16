import "./button";
import { createLabeledIcon } from "./icon";

// Creates a link that looks like a button.
export function createLink(
  icon: string,
  label: string,
  href: string,
  classes: string[] = []
): HTMLAnchorElement {
  const a = document.createElement("a");
  a.classList.add("button");

  for (const className of classes) {
    a.classList.add(className);
  }
  a.href = href;
  a.append(createLabeledIcon(icon, label));
  return a;
}
