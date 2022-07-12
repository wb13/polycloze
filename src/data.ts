// Contains functions for getting data from the server and from localStorage.

import { src } from "./config";
import { Item } from "./item";
import { ItemsSchema, ReviewSchema, SupportedLanguagesSchema } from "./schema";
import { Language } from "./select";

// Local storage stuff

export function getL1 (): string {
    return localStorage.getItem("l1") || "eng";
}

export function getL2 (): string {
    return localStorage.getItem("l2") || "spa";
}

function swapL1L2 () {
    const l1 = getL1();
    const l2 = getL2();
    localStorage.setItem("l1", l2);
    localStorage.setItem("l2", l1);
}

// NOTE Swaps L1 and L2 if code = L2.
// This is needed to make sure the language select form is consistent with
// what's in localStorage.
export function setL1 (code: string) {
    if (code !== getL2()) {
        localStorage.setItem("l1", code);
    } else {
        swapL1L2();
    }
}

// NOTE Swaps L1 and L2 if code = L1.
export function setL2 (code: string) {
    if (code !== getL1()) {
        localStorage.setItem("l2", code);
    } else {
        swapL1L2();
    }
}

function currentCourse (n = 10, x: string[] = []): string {
    let url = `/${getL1()}/${getL2()}?n=${n}`;
    for (const word of x) {
        url += `&x=${word}`;
    }
    return url;
}

// Server stuff

async function fetchJson<T> (url: string | URL, options: RequestInit): Promise<T> {
    if (url instanceof URL) {
        url = url.href;
    }
    const request = new Request(url, options);
    const response = await fetch(request);
    return await response.json();
}

export async function supportedLanguages (): Promise<Language[]> {
    const url = new URL("/options", src);
    const options = { mode: "cors" as RequestMode };
    const json = await fetchJson<SupportedLanguagesSchema>(url, options);
    return json.languages;
}

export async function fetchItems (n = 10, x: string[] = []): Promise<Item[]> {
    const url = new URL(currentCourse(n, x), src);
    const options = { mode: "cors" as RequestMode };
    const json = await fetchJson<ItemsSchema>(url, options);
    return json.items;
}

// Returns response status (success or not).
// Also dispatches a custom event on window.
export async function submitReview (word: string, correct: boolean): Promise<boolean> {
    const url = new URL(currentCourse(), src);
    const options = {
        body: JSON.stringify({
            reviews: [
                { word, correct }
            ]
        }),
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json"
        },
        method: "POST",
        mode: "cors" as RequestMode
    };
    const json = await fetchJson<ReviewSchema>(url, options);
    return json.success;
}
