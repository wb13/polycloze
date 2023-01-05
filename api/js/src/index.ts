import "./index.css";
import "./components";
import { getVoices } from "./tts";

document.documentElement.lang = "en";

if ("serviceWorker" in navigator) {
  // serviceworker has to be at the root.
  navigator.serviceWorker.register("serviceworker.js");
}

// Speeds up subsequent calls.
getVoices();
