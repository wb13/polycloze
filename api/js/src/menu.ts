import "./menu.css";
import { createButton, setButtonLink } from "./button";
import { createIcon, createLabeledIcon } from "./icon";
import { createModal } from "./modal";

function createSignInButton(): HTMLButtonElement {
    const button = setButtonLink(createButton(createLabeledIcon("sign-in", "Sign in")), "/signin");
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}

function createMenu(signedIn: boolean): HTMLDivElement {
    const div = document.createElement("div");
    div.classList.add("menu");
    if (!signedIn) {
        div.append(createSignInButton());
    } else {
        div.append(
            setButtonLink(createButton(createLabeledIcon("house", "Home")), "/"),
            setButtonLink(createButton(createLabeledIcon("brain", "Study")), "/study"),
            setButtonLink(createButton(createLabeledIcon("notebook", "Vocabulary")), "/vocab"),
            setButtonLink(createButton(createLabeledIcon("faders", "Settings")), "/settings"),
            setButtonLink(createButton(createLabeledIcon("sign-out", "Sign out")), "/signout", "POST"),
        );
    }

    const buttons = div.querySelectorAll("button");
    for (let i = 0; i < buttons.length; i++) {
        const button = buttons[i];
        button.classList.add("button-borderless");
        button.classList.add("button-tight");
    }
    return div;
}

function createMenuModal(signedIn: boolean): [HTMLDivElement, () => void] {
    const menu = createMenu(signedIn);
    menu.classList.add("menu-narrow");
    return createModal(menu, {includeHeader: false});
}

function createMenuListButton(signedIn: boolean): HTMLButtonElement {
    const [modal, show] = createMenuModal(signedIn);
    document.body.appendChild(modal);

    const button = createButton(createIcon("list"), show);
    button.classList.add("menu-list-button");
    button.classList.add("button-borderless");
    button.classList.add("button-tight");
    return button;
}

// NOTE requires course-select-button
export function createResponsiveMenu(signedIn: boolean): HTMLDivElement {
    const div = document.createElement("div");
    if (!signedIn) {
        // Just show sign in button.
        div.append(createSignInButton());
        return div;
    }

    const wideMenu = createMenu(signedIn);
    wideMenu.classList.add("menu-wide");
    div.append(
        wideMenu,
        document.createElement("course-select-button"),
        createMenuListButton(signedIn),
    );
    return div;
}
