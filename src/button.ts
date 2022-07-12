import "./button.css";

export function createButton (content: string, onClick?: (event: Event) => void): HTMLButtonElement {
    const button = document.createElement("button");
    button.textContent = content;
    if (onClick) {
        button.addEventListener("click", onClick);
    }
    return button;
}
