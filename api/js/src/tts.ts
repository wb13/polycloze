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

// Async wrapper around `speechSynthesis.getVoices`.
// Returns empty array only if no voices are installed.
// The promise can take a while to resolve on the first call.
// The caller can call it on page load to speed up subsequent calls.
export function getVoices(): Promise<SpeechSynthesisVoice[]> {
  return new Promise((resolve) => {
    const voices = speechSynthesis.getVoices();
    if (voices.length > 0) {
      return resolve(voices);
    }

    const listener = () => {
      resolve(speechSynthesis.getVoices());
      speechSynthesis.removeEventListener("voiceschanged", listener);
    };
    speechSynthesis.addEventListener("voiceschanged", listener);
  });
}

// Returns Map of voices: voiceURI -> voice.
// Includes only voices in the selected language.
async function getVoicesForCurrentLanguage(): Promise<
  Map<string, SpeechSynthesisVoice>
> {
  const lang = getL2().bcp47;
  const map = new Map();
  for (const voice of await getVoices()) {
    if (compareVoice(lang, voice)) {
      map.set(voice.voiceURI, voice);
    }
  }
  return map;
}

// Usage: call `init()` after constructor.
export class TTS {
  voices: Map<string, SpeechSynthesisVoice>;
  private initialized = false;
  private startCallbacks: Array<() => void> = [];
  private endCallbacks: Array<() => void> = [];

  constructor() {
    this.voices = new Map();
  }

  async init() {
    this.voices = await getVoicesForCurrentLanguage();
    this.initialized = true;
  }

  // Speaks text using the preferred voice if TTS is enabled.
  speak(text: string) {
    // Make sure TTS was initialized.
    if (!this.initialized) {
      throw new Error("TTS object was not initialized");
    }

    // Check if TTS is enabled.
    if (!isEnabledTTS()) {
      return;
    }

    // Create utterance.
    const utterance = new SpeechSynthesisUtterance(text);
    const voice = this.voices.get(getPreferredVoice() || "");
    if (voice != null) {
      utterance.voice = voice;
    } else if (this.voices.size > 0) {
      utterance.voice = this.voices.values().next().value;
    } else {
      utterance.lang = getL2().bcp47;
    }

    // Set utterance callbacks.
    utterance.addEventListener("start", () => {
      for (const callback of this.startCallbacks) {
        callback();
      }
    });
    utterance.addEventListener("end", () => {
      for (const callback of this.endCallbacks) {
        callback();
      }
    });

    // Stop currently playing utterance and play the next one.
    this.stop();
    speechSynthesis.speak(utterance);
  }

  stop() {
    speechSynthesis.cancel();
  }

  // Add an utterance callback function.
  onStart(callback: () => void) {
    this.startCallbacks.push(callback);
  }

  // Add an utterance callback function.
  onEnd(callback: () => void) {
    this.endCallbacks.push(callback);
  }
}

// Returns URI of preferred voice for current language.
function getPreferredVoice(): string | null {
  const lang = getL2();
  return localStorage.getItem(`tts.${lang.code}.voiceURI`);
}

// Sets voice as preferred voice for current language.
// Assumes the voice is in the correct language.
function setPreferredVoice(voiceURI: string) {
  const lang = getL2();
  localStorage.setItem(`tts.${lang.code}.voiceURI`, voiceURI);
}

function createVoiceSelect(tts: TTS): HTMLSelectElement {
  const preferred = getPreferredVoice();

  const select = document.createElement("select");
  if (!isEnabledTTS()) {
    select.disabled = true;
  }

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

  // If there are no available voices, say so.
  if (select.childElementCount <= 0) {
    const option = document.createElement("option");
    option.textContent = "(No voices available.)";
    select.appendChild(option);
  }

  select.addEventListener("change", () => setPreferredVoice(select.value));
  return select;
}

function createVoicePlayButton(tts: TTS): HTMLButtonElement {
  let playing = false;
  let icon = createIcon("speaker-high");
  tts.onStart(() => {
    const replacement = createIcon("stop");
    icon.replaceWith(replacement);
    icon = replacement;
    playing = true;
  });
  tts.onEnd(() => {
    const replacement = createIcon("speaker-high");
    icon.replaceWith(replacement);
    icon = replacement;
    playing = false;
  });

  const button = createButton(icon, async () => {
    if (playing) {
      tts.stop();
      return;
    }
    const sentences = await fetchSentences();
    if (sentences.length === 0) {
      return;
    }
    tts.speak(sentences[0].text);
  });
  button.type = "button";
  button.classList.add("button-tight");

  if (!isEnabledTTS()) {
    button.disabled = true;
  }
  return button;
}

// Returns a div element, and a function for enabling/disabling the demo.
function createVoiceDemo(
  tts: TTS
): [HTMLDivElement, (checked: boolean) => void] {
  const div = document.createElement("div");
  div.classList.add("tts-demo");

  const button = createVoicePlayButton(tts);
  const select = createVoiceSelect(tts);
  div.append(button, select);

  const hook = (checked: boolean) => {
    if (checked) {
      button.disabled = false;
      select.disabled = false;
    } else {
      button.disabled = true;
      select.disabled = true;
    }
  };
  return [div, hook];
}

// Enables TTS in the selected language.
function enableTTS() {
  const lang = getL2();

  // Default value is falsy if not set, so we use the `disabled` field.
  localStorage.setItem(`tts.${lang.code}.disabled`, "false");
}

// Disables TTS in the selected language.
function disableTTS() {
  const lang = getL2();
  localStorage.setItem(`tts.${lang.code}.disabled`, "true");
}

// Returns whether or not TTS is disabled for the selected language.
export function isEnabledTTS(): boolean {
  // Local storage stores `disabled` instead of `enabled`, because TTS is
  // enabled by default. So when the item in the local storage isn't set,
  // it is enabled.
  const lang = getL2();
  return localStorage.getItem(`tts.${lang.code}.disabled`) === "true"
    ? false
    : true;
}

// Takes a callback function that gets called when the checkbox gets clicked.
// The callback function takes a boolean (checked or not).
function createVoiceCheckbox(
  disabled: boolean,
  callback: (checked: boolean) => void
): HTMLDivElement {
  const div = document.createElement("div");
  div.innerHTML = `
        <input type="checkbox" id="enable-tts" name="enable-tts">
        <label for="enable-tts">Enable text-to-speech</label>
    `;

  const input = div.querySelector("input") as HTMLInputElement;
  if (disabled) {
    input.disabled = true;
  }

  if (isEnabledTTS()) {
    input.checked = true;
  }
  input.addEventListener("click", () => {
    callback(input.checked);
    if (input.checked) {
      enableTTS();
    } else {
      disableTTS();
    }
  });
  return div;
}

export function createVoiceSettingsSection(tts: TTS): HTMLFormElement {
  const form = document.createElement("form");
  form.classList.add("signin");

  const disabled = tts.voices.size === 0;
  if (disabled) {
    disableTTS();
  }

  const [demo, hook] = createVoiceDemo(tts);
  form.append(createVoiceCheckbox(disabled, hook), demo);
  return form;
}
