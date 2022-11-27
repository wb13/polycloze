import "./menu.css";
import { createButton, setButtonLink } from "./button";
import { createIcon, createLabeledIcon } from "./icon";
import { createModal } from "./modal";

function createMenu(signedIn: boolean): HTMLDivElement {
    const div = document.createElement("div");
    div.classList.add("menu");
    if (!signedIn) {
        div.append(
            setButtonLink(createButton(createLabeledIcon("sign-in", "Sign in")), "/signin"),
        );
        return div;
    }
    div.append(
        setButtonLink(createButton(createLabeledIcon("house", "Home")), "/"),
        setButtonLink(createButton(createLabeledIcon("brain", "Study")), "/study"),
        setButtonLink(createButton(createLabeledIcon("notebook", "Vocabulary")), "/vocab"),
        setButtonLink(createButton(createLabeledIcon("faders", "Settings")), "/settings"),
        setButtonLink(createButton(createLabeledIcon("sign-out", "Sign out")), "/signout"),
    );

    const buttons = div.querySelectorAll("button");
    for (let i = 0; i < buttons.length; i++) {
        const button = buttons[i];
        button.classList.add("button-borderless");
        button.classList.add("button-tight");
    }
    return div;
}

function createMenuModal(signedIn: boolean): [HTMLDivElement, () => void] {
    return createModal(createMenu(signedIn), {includeHeader: false});
}

export function createMenuListButton(signedIn: boolean): HTMLButtonElement {
    const [modal, show] = createMenuModal(signedIn);
    document.body.appendChild(modal);

    const button = createButton(createIcon("list"), show);
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}
