import { getL2 } from "./language";

function fixLanguage(lang: string): string {
    // Speech synthesis voices sometimes don't use BCP47 language codes.
    if (lang === "tl") {
        const choices = new Set(["tl", "tgl", "fil", "fil-PH"]);
        for (const voice of speechSynthesis.getVoices()) {
            if (choices.has(voice.lang)) {
                return voice.lang;
            }
        }
    }
    return lang;
}

export function speak(text: string) {
    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = fixLanguage(getL2().bcp47);
    speechSynthesis.speak(utterance);
}
