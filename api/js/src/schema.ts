// JSON schemas used by server.

import { Item } from "./item";
import { Language } from "./language";

export type ItemsSchema = {
    items: Item[];
};

export type ReviewSchema = {
    success: boolean;
    frequencyClass: number;    // describes student's level
};

export type CourseStats = {
  seen?: number
  total?: number
  learned?: number
  reviewed?: number
  correct?: number
};

export type Course = {
    l1: Language;
    l2: Language;
    stats?: CourseStats;
};

export type CoursesSchema = {
    courses: Course[];
};

export type VocabularyItem = {
  word: string;
  reviewed: string;
  due: string;
  strength: number;
};

// from /<l1>/<l2>/vocab
export type VocabularySchema = {
  results: VocabularyItem[];
};
