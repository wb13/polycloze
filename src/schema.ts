// JSON schemas used by server.

import { Item } from "./item";
import { Language } from "./select";

export type ItemsSchema = {
    items: Item[];
};

export type ReviewSchema = {
    success: boolean;
    frequencyClass: number;    // describes student's level
};

export type SupportedLanguagesSchema = {
    languages: Language[];
};
