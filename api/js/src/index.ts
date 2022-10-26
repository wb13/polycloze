import "./index.css";
import "./components.ts";

document.documentElement.lang = "en";

if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("dist/serviceworker.js");
}
