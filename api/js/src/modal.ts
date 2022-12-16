import "./modal.css";
import { createButton } from "./button";
import { createIcon } from "./icon";

function createModalCloseButton(
  onClick: (event: Event) => void
): HTMLButtonElement {
  const button = createButton(createIcon("x"), onClick);
  button.classList.add("button-borderless");
  button.classList.add("button-tight");
  return button;
}

function createModalHeader(hide: () => void): HTMLDivElement {
  const div = document.createElement("div");
  div.classList.add("modal-header");
  div.appendChild(createModalCloseButton(hide));
  return div;
}

type CreateModalDialogOptions = {
  includeHeader?: boolean;
};

function defaultCreateModalDialogOptions(): CreateModalDialogOptions {
  return { includeHeader: true };
}

function createModalDialog(
  body: string | Element,
  hide: () => void,
  options: CreateModalDialogOptions = {}
): HTMLDivElement {
  const { includeHeader } = {
    ...defaultCreateModalDialogOptions(),
    ...options,
  };
  const div = document.createElement("div");
  div.classList.add("modal-dialog");
  if (includeHeader) {
    div.appendChild(createModalHeader(hide));
  }
  div.append(body);
  return div;
}

function createModalBackground(hide: () => void): HTMLDivElement {
  const div = document.createElement("div");
  div.classList.add("modal-background");
  div.addEventListener("click", hide);
  return div;
}

const styleReset: Array<() => void> = [];

// Used for showing modal elements only.
function showModalElement(element: HTMLElement) {
  element.classList.remove("modal-hidden");

  // Save original style values.
  const style = {
    overflow: document.body.style.overflow,
    paddingRight: document.body.style.paddingRight,
    width: document.body.style.width,
  };

  // Prevent whole page from shifting to the right when scroll bar disappears.
  // Idea: widen the document body to fill the space left behind by the scrollbar.
  const htmlWidth = document.documentElement.clientWidth;
  const windowWidth = window.innerWidth;
  const scrollbarWidth = windowWidth - htmlWidth;
  if (scrollbarWidth > 0) {
    document.body.style.paddingRight = Math.floor(scrollbarWidth) + "px";
    document.body.style.width = `calc(100vw - ${Math.floor(scrollbarWidth)}px)`;
  }

  // Prevent scrolling.
  document.body.style.overflow = "hidden";

  // Function that undoes CSS changes in this function.
  const reset = () => {
    const { overflow, paddingRight, width } = style;
    document.body.style.overflow = overflow;
    document.body.style.paddingRight = paddingRight;
    document.body.style.width = width;
  };
  styleReset.push(reset);
}

// Used for hiding modal elements only.
function hideModalElement(element: HTMLElement) {
  element.classList.add("modal-hidden");
  document.body.style.overflow = "";
  const reset = styleReset.pop();
  if (reset != null) {
    reset();
  }
}

type CreateModalOptions = {
  includeHeader?: boolean;
};

function defaultCreateModalOptions(): CreateModalOptions {
  return { includeHeader: true };
}

// Creates modal.
// Returns a div element and a show-function.
// Modal should be inserted into document.body.
export function createModal(
  body: string | Element,
  options: CreateModalOptions = {}
): [HTMLDivElement, () => void] {
  const { includeHeader } = { ...defaultCreateModalOptions(), ...options };

  const div = document.createElement("div");
  div.classList.add("modal");

  const show = () => showModalElement(div);
  const hide = () => hideModalElement(div);
  hide();

  div.append(
    createModalBackground(hide),
    createModalDialog(body, hide, { includeHeader })
  );
  return [div, show];
}
