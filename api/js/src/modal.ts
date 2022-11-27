import "./modal.css";
import { createButton } from "./button";
import { createIcon } from "./icon";

function createModalCloseButton(onClick: (event: Event) => void): HTMLButtonElement {
    const button = createButton(createIcon("x"), onClick);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}

function createModalHeader(hide: () => void): HTMLDivElement {
    // TODO content
    const div = document.createElement("div");
    div.classList.add("modal-header");
    div.appendChild(createModalCloseButton(hide));
    return div;
}

type CreateModalDialogOptions = {
    includeHeader?: boolean;
};

function defaultCreateModalDialogOptions(): CreateModalDialogOptions {
    return {includeHeader: true};
}

function createModalDialog(body: string | Element, hide: () => void, options: CreateModalDialogOptions = {}): HTMLDivElement {
    const {includeHeader} = {
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

// Used for showing modal elements only.
function showModalElement(element: HTMLElement) {
    element.style.display = "";
    document.body.style.overscrollBehavior = "contain";
    document.body.style.height = "100vh";
    document.body.style.overflowY = "hidden";
    document.body.style.width = "100vw";
    document.body.style.overflowX = "hidden";
    document.body.style.position = "fixed";
}

// Used for hiding modal elements only.
function hideModalElement(element: HTMLElement) {
    element.style.display = "none";
    document.body.style.overscrollBehavior = "";
    document.body.style.height = "";
    document.body.style.overflowY = "";
    document.body.style.width = "";
    document.body.style.overflowX = "";
    document.body.style.position = "";
}

type CreateModalOptions = {
    includeHeader?: boolean;
};

function defaultCreateModalOptions(): CreateModalOptions {
    return {includeHeader: true};
}

// Creates modal.
// Returns a div element and a show-function.
// Modal should be inserted into document.body.
export function createModal(body: string | Element, options: CreateModalOptions = {}): [HTMLDivElement, () => void] {
    const { includeHeader } = {...defaultCreateModalOptions(), ...options};

    const div = document.createElement("div");
    div.classList.add("modal");

    const show = () => showModalElement(div);
    const hide = () => hideModalElement(div);
    hide();

    div.append(
        createModalBackground(hide),
        createModalDialog(body, hide, { includeHeader }),
    );
    return [div, show];
}
