// Wrappers for api calls.

import { csrf } from "./csrf";
import { day, endOfDay } from "./datetime";
import { Difficulty } from "./difficulty";
import { Item } from "./item";
import { getL1, getL2 } from "./language";
import { fetchJson, resolve, submitJson } from "./request";
import {
  ActivitySchema,
  ActivitySummary,
  Course,
  CoursesSchema,
  DataPoint,
  EstimatedLevelSchema,
  FlashcardsResponse,
  Language,
  LanguagesSchema,
  RandomSentence,
  RandomSentencesSchema,
  ReviewResult,
  SetCourseResponse,
  Word,
  VocabularySchema,
  VocabularySizeSchema,
} from "./schema";

type FetchVocabularyOptions = {
  // Path params
  l1?: string; // L1 code
  l2?: string; // L2 code

  // Search params
  limit?: number; // Max number of items to fetch
  after?: string; // Last item to exclude from query
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

export async function fetchVocabulary(
  options: FetchVocabularyOptions = {}
): Promise<Word[]> {
  const { l1, l2, limit, after, sortBy } = {
    ...defaultFetchVocabularyOptions(),
    ...options,
  };
  const url = resolve(`/${l1}/${l2}/vocab`);
  setParams(url, { after, limit, sortBy });

  const json = await fetchJson<VocabularySchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.words || [];
}

type FetchActivityOptions = {
  l1?: string;
  l2?: string;
  from?: Date;
  to?: Date;
  step?: number;
};

function defaultFetchActivityOptions(): FetchActivityOptions {
  const to = endOfDay();
  const from = new Date(to.valueOf() - 7 * day);
  return {
    l1: getL1().code,
    l2: getL2().code,
    from,
    to,
    step: 86400, // 1 day
  };
}

// Fetches student's recent activity.
export async function fetchActivity(
  options: FetchActivityOptions = {}
): Promise<ActivitySummary[]> {
  options = { ...defaultFetchActivityOptions(), ...options };
  const { l1, l2 } = options;
  const url = resolve(`/api/stats/activity/${l1}/${l2}`);
  setParams(url, {
    from: options.from ? options.from.getTime() / 1000 : undefined,
    to: options.to ? options.to.getTime() / 1000 : undefined,
    step: options.step || undefined,
  });
  const json = await fetchJson<ActivitySchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.activity.map((s) => {
    return {
      ...s,
      from: new Date(s.from),
      to: new Date(s.to),
    };
  });
}

type FetchHistoricalDataOptions = {
  l1?: string;
  l2?: string;
  from?: Date;
  to?: Date;
  step?: number;
};

function defaultFetchHistoricalDataOptions(): FetchHistoricalDataOptions {
  const to = endOfDay();
  const from = new Date(to.valueOf() - 7 * day);
  return {
    l1: getL1().code,
    l2: getL2().code,
    from,
    to,
    step: 86400, // 1 day
  };
}

// Fetches student's vocab size over time.
export async function fetchVocabularySize(
  options: FetchHistoricalDataOptions = {}
): Promise<DataPoint[]> {
  options = { ...defaultFetchHistoricalDataOptions(), ...options };
  const { l1, l2 } = options;
  const url = resolve(`/api/stats/vocab/${l1}/${l2}`);
  setParams(url, {
    from: options.from ? options.from.getTime() / 1000 : undefined,
    to: options.to ? options.to.getTime() / 1000 : undefined,
    step: options.step || undefined,
  });
  const json = await fetchJson<VocabularySizeSchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.vocabSize.map((p) => {
    return {
      time: new Date(p.time),
      value: p.value,
    };
  });
}

// Fetches student's estimated level over time.
export async function fetchEstimatedLevel(
  options: FetchHistoricalDataOptions = {}
): Promise<DataPoint[]> {
  options = { ...defaultFetchHistoricalDataOptions(), ...options };
  const { l1, l2 } = options;
  const url = resolve(`/api/stats/estimate/${l1}/${l2}`);
  setParams(url, {
    from: options.from ? options.from.getTime() / 1000 : undefined,
    to: options.to ? options.to.getTime() / 1000 : undefined,
    step: options.step || undefined,
  });
  const json = await fetchJson<EstimatedLevelSchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.estimatedLevel.map((p) => {
    return {
      time: new Date(p.time),
      value: p.value,
    };
  });
}

export async function fetchCourses(): Promise<Course[]> {
  const url = resolve("/api/courses");
  setParams(url, { t: "20221114" });
  const json = await fetchJson<CoursesSchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.courses;
}

// Fetches list of supported languages (L1).
export async function fetchLanguages(): Promise<Language[]> {
  const url = resolve("/api/languages");
  setParams(url, { t: "20221114" });
  const json = await fetchJson<LanguagesSchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.languages;
}

type FetchFlashcardsOptions = {
  // Path params
  l1?: string; // L1 code
  l2?: string; // L2 code

  // Body params
  limit?: number; // Max number of flashcards to fetch
  exclude?: string[]; // Words to exclude in flashcards
  reviews?: ReviewResult[];
  difficulty?: Difficulty;
};

function defaultFetchFlashcardsOptions(): FetchFlashcardsOptions {
  return {
    l1: getL1().code,
    l2: getL2().code,
    limit: 10,
    exclude: [],
    reviews: [],
  };
}

// Returns a copy of the review result containing only the necessary fields.
function minimizeReviewResult(review: ReviewResult): ReviewResult {
  const { word, correct, timestamp } = review;
  return { word, correct, timestamp };
}

export function fetchFlashcards(
  options: FetchFlashcardsOptions = {}
): Promise<FlashcardsResponse> {
  options = { ...defaultFetchFlashcardsOptions(), ...options };
  const { l1, l2 } = options;
  const url = resolve(`/api/flashcards/${l1}/${l2}`);
  const data = {
    limit: options.limit,
    exclude: options.exclude,
    reviews:
      options.reviews != null
        ? options.reviews.map(minimizeReviewResult)
        : undefined,
    difficulty: options.difficulty,
    timestamp: Math.floor(Date.now() / 1000),
  };
  return submitJson<FlashcardsResponse>(url, data);
}

// Sends review results to the server.
// It uses the `sendBeacon` function to make sure the data gets sent to the
// server.
// This can be safely used inside a `visibilitychange` listener to upload
// review results before the browser gets closed.
export function sendReviewResults(
  reviews: ReviewResult[],
  difficulty: Difficulty
) {
  const l1 = getL1().code;
  const l2 = getL2().code;
  const url = resolve(`/api/flashcards/${l1}/${l2}`);
  const data = {
    limit: 0,
    reviews,
    difficulty,
    csrfToken: csrf(),
    timestamp: Math.floor(Date.now() / 1000),
  };
  const blob = new Blob([JSON.stringify(data)], {
    type: "application/json",
  });
  navigator.sendBeacon(url, blob);
}

type FetchSentencesOptions = {
  l1?: string;
  l2?: string;
  limit?: number;
};

function defaultFetchSentencesOptions(): FetchSentencesOptions {
  return {
    l1: getL1().code,
    l2: getL2().code,
    limit: 1,
  };
}

export async function fetchSentences(
  options: FetchSentencesOptions = {}
): Promise<RandomSentence[]> {
  const { l1, l2, limit } = { ...defaultFetchSentencesOptions(), ...options };
  if (l1 == null || l2 == null) {
    throw new Error("l1 and l2 required");
  }

  const url = resolve("/api/sentences");
  setParams(url, { l1, l2, limit });

  const json = await fetchJson<RandomSentencesSchema>(url, {
    mode: "cors" as RequestMode,
  });
  return json.sentences;
}

export async function setActiveCourse(
  l1: string,
  l2: string
): Promise<boolean> {
  const url = resolve("/api/actions/set-course");
  const data = {
    l1Code: l1,
    l2Code: l2,
  };
  const resp = await submitJson<SetCourseResponse>(url, data);
  return resp.ok;
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
