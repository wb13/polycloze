// Contains functions for getting data from the server and from localStorage.

import { Item } from "./item";
import { getL1, getL2 } from "./language";
import { Course, ItemsSchema, ReviewSchema, CoursesSchema } from "./schema";

// Location of server
const src = findServer();

function currentCourse(n = 10, x: string[] = []): string {
    const l1 = getL1().code;
    const l2 = getL2().code;
    let url = `/${l1}/${l2}?n=${n}`;
    for (const word of x) {
        url += `&x=${word}`;
    }
    return url;
}

// Server stuff

async function fetchJson<T>(url: string | URL, options: RequestInit): Promise<T> {
    if (url instanceof URL) {
        url = url.href;
    }
    const request = new Request(url, options);
    const response = await fetch(request);
    return await response.json();
}

export type AvailableCoursesOptions = {
    l1?: string;
    l2?: string;
    stats?: boolean;
};

export async function availableCourses(params: AvailableCoursesOptions = {}): Promise<Course[]> {
    const url = new URL("/courses", src);

    if (params.l1 != null && params.l1.length > 0) {
        url.searchParams.set("l1", params.l1);
    }
    if (params.l2 != null && params.l2.length > 0) {
        url.searchParams.set("l2", params.l2);
    }
    if (params.stats) {
        url.searchParams.set("stats", "true");
    }

    const options = { mode: "cors" as RequestMode };
    const json = await fetchJson<CoursesSchema>(url, options);
    return json.courses;
}

export async function fetchItems(n = 10, x: string[] = []): Promise<Item[]> {
    const url = new URL(currentCourse(n, x), src);
    const options = { mode: "cors" as RequestMode };
    const json = await fetchJson<ItemsSchema>(url, options);
    return json.items;
}

// Returns response status (success or not).
// Also dispatches a custom event on window.
export async function submitReview(word: string, correct: boolean): Promise<ReviewSchema> {
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
    return await fetchJson<ReviewSchema>(url, options);
}

function findServer(): string {
    const url = new URL(location.href);
    if (document.currentScript == null) {
        return url.origin;
    }

    const { origin, port } = document.currentScript.dataset;
    if (origin != null) {
        return origin;
    }
    if (port != null) {
        url.port = port;
    }
    return url.origin;
}
