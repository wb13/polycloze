import "./button.css";
import { createCSRFTokenInput } from "./csrf";

type ButtonContent = string | Element | DocumentFragment;

export function setButton(button: HTMLButtonElement, content: ButtonContent, onClick?: (event: Event) => void) {
    button.append(content);
    if (onClick) {
        button.addEventListener("click", onClick);
    }
}

export function createButton(content: ButtonContent, onClick?: (event: Event) => void): HTMLButtonElement {
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

export function setButtonLink(button: HTMLButtonElement, href: string, method = "GET"): HTMLButtonElement {
    button.addEventListener("click", () => {
        if (method === "POST") {
            followLinkPost(href);
        } else {
            window.location.pathname = href;
        }
    });
    return button;
}
