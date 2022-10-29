import "./button.css";

export function setButton(button: HTMLButtonElement, content: string, onClick?: (event: Event) => void) {
    button.textContent = content;
    if (onClick) {
        button.addEventListener("click", onClick);
    }
}

export function createButton(content: string, onClick?: (event: Event) => void): HTMLButtonElement {
    const button = document.createElement("button");
    setButton(button, content, onClick);
    return button;
}
