// JSON schemas used by server.

import { Difficulty } from "./difficulty";
import { Item } from "./item";

export type ItemsSchema = {
  items: Item[];
};

export type FlashcardsResponse = {
  items: Item[];
  difficulty: Difficulty;
};

export type SetCourseRequest = {
  l1Code: string;
  l2Code: string;
};

export type SetCourseResponse = {
  ok: boolean;
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

export type ActivitySummary = {
  from: Date;
  to: Date;

  unimproved: number;
  learned: number;
  forgotten: number;
  crammed: number;
  strengthened: number;
};

// Same as ActivitySummary, but with unparsed timestamps.
export type ActivitySummarySchema = {
  from: string;
  to: string;

  unimproved: number;
  learned: number;
  forgotten: number;
  crammed: number;
  strengthened: number;
};

// from /api/stats/activity/<l1>/<l2>?from=<from>&to=<to>&step=<step>
export type ActivitySchema = {
  activity: ActivitySummarySchema[];
};

// Not to be confused with sentence.Sentence.
export type RandomSentence = {
  id: number;
  tatoebaID?: number;
  text: string;
};

export type RandomSentencesSchema = {
  sentences: RandomSentence[];
};

export type DataPoint = {
  time: Date;
  value: number;
};

// Same as DataPoint, but with unparsed timestamp.
export type DataPointSchema = {
  time: string;
  value: number;
};

// from /api/stats/vocab/<l1>/<l2>?from=<from>&to=<to>&step=<step>
export type VocabularySizeSchema = {
  vocabSize: DataPointSchema[];
};

// from /api/stats/estimate/<l1>/<l2>
export type EstimatedLevelSchema = {
  estimatedLevel: DataPointSchema[];
};

export type ReviewResult = {
  word: string;
  correct: boolean;
  timestamp: number;

  // This field doesn't need to be sent to the server.
  new?: boolean;
};

export type UploadCSVFileResponse = {
  message: string;
  success: boolean;
};
