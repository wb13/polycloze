import "./tts.css";
import { fetchSentences } from "./api";
import { createButton } from "./button";
import { createIcon } from "./icon";
import { getL1, getL2 } from "./language";

function rawCompareVoice(lang: string, voice: SpeechSynthesisVoice): boolean {
    return voice.lang === lang || voice.lang.startsWith(lang + "-");
}

// Checks if voice is in the given language.
function compareVoice(lang: string, voice: SpeechSynthesisVoice): boolean {
    // Tagalog correction.
    if (lang === "tl") {
        const langs = ["tl", "fil-PH", "tgl", "fil"];
        for (const code of langs) {
            if (rawCompareVoice(code, voice)) {
                return true;
            }
        }
    }
    return rawCompareVoice(lang, voice);
}

// Returns Map of voices: voiceURI -> voice.
// Includes only voices in the selected language.
function getVoices(): Map<string, SpeechSynthesisVoice> {
    const lang = getL2().bcp47;
    const map = new Map();
    for (const voice of speechSynthesis.getVoices()) {
        if (compareVoice(lang, voice)) {
            map.set(voice.voiceURI, voice);
        }
    }
    return map;
}

export class TTS {
    voices: Map<string, SpeechSynthesisVoice>;

    constructor() {
        this.voices = getVoices();
    }

    speak(text: string) {
        const utterance = new SpeechSynthesisUtterance(text);

        const voice = this.voices.get(getPreferredVoice() || "");
        if (voice != null) {
            utterance.voice = voice;
        } else if (this.voices.size > 0) {
            utterance.voice = this.voices.values().next().value;
        } else {
            utterance.lang = getL2().bcp47;
        }
        speechSynthesis.speak(utterance);
    }
}

// Returns URI of preferred voice for current language.
function getPreferredVoice(): string | null {
    const lang = getL2();
    return localStorage.getItem(`voice.${lang.bcp47}`);
}

// Sets voice as preferred voice for current language.
// Assumes the voice is in the correct language.
function setPreferredVoice(voiceURI: string) {
    const lang = getL2();
    localStorage.setItem(`voice.${lang.bcp47}`, voiceURI);
}

function createVoiceSelect(tts: TTS): HTMLSelectElement {
    const preferred = getPreferredVoice();

    const select = document.createElement("select");

    for (const voice of tts.voices.values()) {
        const option = document.createElement("option");
        option.value = voice.voiceURI;
        option.textContent = voice.name;
        if (preferred === voice.voiceURI) {
            option.selected = true;
        }
        select.appendChild(option);
    }
    // TODO what if none selected?

    select.addEventListener("change", () => setPreferredVoice(select.value));
    return select;
}

function createVoicePlayButton(tts: TTS): HTMLButtonElement {
    const button = createButton(createIcon("speaker-high"), async() => {
        // Voice demo.
        const sentences = await fetchSentences();
        if (sentences.length === 0) {
            return;
        }
        tts.speak(sentences[0].text);
    });
    button.type = "button";
    button.classList.add("button-tight");
    return button;
}

function createVoiceDemo(tts: TTS): HTMLDivElement {
    const div = document.createElement("div");
    div.classList.add("tts-demo");
    div.append(createVoicePlayButton(tts), createVoiceSelect(tts));
    return div;
}

export function createVoiceSettingsSection(): HTMLFormElement {
    const tts = new TTS();

    const form = document.createElement("form");
    form.classList.add("signin");

    const h2 = document.createElement("h2");
    h2.textContent = `${getL2().name} from ${getL1().name} settings`;

    form.append(
        h2,
        createVoiceDemo(tts),
    );
    return form;
}
