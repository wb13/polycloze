// JSON schemas used by server.

import { Item } from "./item";

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

export type Language = {
  code: string;
  name: string;
  bcp47: string;
};

export type LanguagesSchema = {
  languages: Language[];
};

export type Course = {
    l1: Language;
    l2: Language;
    stats?: CourseStats;
};

export type CoursesSchema = {
    courses: Course[];
};

export type Word = {
  word: string;
  learned: string;
  reviewed: string;
  due: string;
  strength: number;
};

// from /<l1>/<l2>/vocab
export type VocabularySchema = {
  words: Word[];
};

export type Activity = {
  forgotten: number;
  unimproved: number;
  crammed: number;
  learned: number;
  strengthened: number;
};

// from /<l1>/<l2>/activity
export type ActivityHistory = {
  activities: Activity[]; // up to one year of activities
  aggregates: Activity;   // for > 1 year old
};
