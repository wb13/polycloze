import "./button.css";
import { createCSRFTokenInput } from "./csrf";

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

function followLinkPost(url: string) {
    const form = document.createElement("form");
    form.action = url;
    form.method = "POST";
    form.style.display = "none";
    form.appendChild(createCSRFTokenInput());

    document.body.appendChild(form);

    form.submit();
}

export function setButtonLink(button: HTMLButtonElement, href: string, method = "GET") {
    button.addEventListener("click", () => {
        if (method === "POST") {
            followLinkPost(href);
        } else {
            window.location.pathname = href;
        }
    });
}
