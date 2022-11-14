// Wrappers for api calls.

import { Item } from "./item";
import { getL1, getL2 } from "./language";
import { fetchJson, resolve, submitJson } from "./request";
import {
    Course,
    CoursesSchema,
    ItemsSchema,
    Language,
    LanguagesSchema,
    ReviewSchema,
    Word,
    VocabularySchema,
} from "./schema";

type FetchVocabularyOptions = {
    // Path params
    l1?: string;    // L1 code
    l2?: string;    // L2 code

    // Search params
    limit?: number;  // Max number of items to fetch
    after?: string;  // Last item to exclude from query
    sortBy?: "word" | "reviewed" | "due" | "strength";
};

function defaultFetchVocabularyOptions(): FetchVocabularyOptions {
    return {
        l1: getL1().code,
        l2: getL2().code,
        limit: 50,
        after: "",
        sortBy: "word",
    };
}

export async function fetchVocabulary(options: FetchVocabularyOptions = {}): Promise<Word[]> {
    const { l1, l2, limit, after, sortBy } = {...defaultFetchVocabularyOptions(), ...options};
    const url = resolve(`/${l1}/${l2}/vocab`);
    setParams(url, { after, limit, sortBy });

    const json = await fetchJson<VocabularySchema>(url, {
        mode: "cors" as RequestMode,
    });
    return json.words || [];
}

// Fetches list of supported languages (L1).
export async function fetchLanguages(): Promise<Language[]> {
    const url = resolve("/share/languages.json");
    setParams(url, { t: "20221114" });
    const json = await fetchJson<LanguagesSchema>(url, {
        mode: "cors" as RequestMode,
    });
    return json.languages;
}

type FetchCoursesOptions = {
    l1?: string;        // L1 code
    l2?: string;        // L2 code
    stats?: boolean;    // Whether to include stats or not
};

function defaultFetchCoursesOptions(): FetchCoursesOptions {
    return { stats: true };
}

export async function fetchCourses(options: FetchCoursesOptions = {}): Promise<Course[]> {
    const params = {...defaultFetchCoursesOptions(), ...options};
    const url = resolve("/courses");
    setParams(url, params);

    const json = await fetchJson<CoursesSchema>(url, {
        mode: "cors" as RequestMode,
    });
    return json.courses;
}

type FetchItemsOptions = {
    // Path params
    l1?: string;    // L1 code
    l2?: string;    // L2 code

    // Search params
    n?: number;      // Max number of items to fetch
    x?: string[];    // Words to exclude
};

function defaultFetchItemsOptions(): FetchItemsOptions {
    return {
        l1: getL1().code,
        l2: getL2().code,
        n: 10,
        x: [],
    };
}

export async function fetchItems(options: FetchItemsOptions = {}): Promise<Item[]> {
    const { l1, l2, n, x } = {...defaultFetchItemsOptions(), ...options};
    const url = resolve(`/${l1}/${l2}`);
    setParams(url, { n, x });

    const json = await fetchJson<ItemsSchema>(url, {
        mode: "cors" as RequestMode,
    });
    return json.items;
}

type Params = {
    [name: string]: unknown;
};

function setParams(url: URL, params: Params) {
    for (const name of Object.getOwnPropertyNames(params)) {
        const value = params[name];
        if (value === undefined) {
            continue;
        }
        if (value instanceof Array) {
            for (const item of value) {
                url.searchParams.append(name, item);
            }
            continue;
        }
        url.searchParams.set(name, String(value));
    }
}

export function submitReview(word: string, correct: boolean): Promise<ReviewSchema> {
    const l1 = getL1().code;
    const l2 = getL2().code;

    const url = resolve(`/${l1}/${l2}`);
    const data = {
        reviews: [
            { word, correct },
        ],
    };
    return submitJson<ReviewSchema>(url, data);
}
