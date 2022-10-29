import "./index.css";
import "./components";

document.documentElement.lang = "en";

if ("serviceWorker" in navigator) {
    // serviceworker has to be at the root.
    navigator.serviceWorker.register("serviceworker.js");
}
