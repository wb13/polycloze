import "./upload.css";
import { uploadCSVFile } from "./api";
import { createButton } from "./button";
import { createLabeledIcon } from "./icon";
import { resolve } from "./request";

// Checks if the file the user wants to upload is too big.
function isTooBig(file: File): boolean {
  // Returns true if file size > 8MB.
  return file.size > 8 * 1024 * 1024;
}

// Creates a hidden file input element for uploading CSV files.
function createHiddenFileInput(
  name: string,
  onError: (message: string) => void
): HTMLInputElement {
  const input = document.createElement("input");
  input.name = name;
  input.type = "file";
  input.accept = ".csv";
  input.required = true;
  input.hidden = true;

  input.addEventListener("change", () => {
    const files = input.files || [];
    const file = files[0];
    uploadFile(name, files[0], onError);
  });
  return input;
}

function createDragAndDropText(): HTMLSpanElement {
  const span = document.createElement("span");
  span.style.fontSize = "1.25em";
  span.style.textAlign = "center";
  span.textContent = "Drag and drop CSV file to upload or";
  return span;
}

function createBrowseFilesButton(
  onClick?: (event: Event) => void
): HTMLButtonElement {
  const label = createLabeledIcon("magnifying-glass", "Browse files");
  const button = createButton(label, onClick);
  button.classList.add("button-tight");
  return button;
}

function createFileBrowserBody(
  onClick?: (event: Event) => void
): HTMLDivElement {
  const div = document.createElement("div");
  div.classList.add("file-browser-body");
  div.appendChild(createDragAndDropText());
  div.appendChild(createBrowseFilesButton(onClick));
  return div;
}

function createErrorMessage(message: string): HTMLDivElement {
  const div = document.createElement("div");
  div.classList.add("incorrect");
  div.appendChild(createLabeledIcon("warning-circle", message));
  return div;
}

export function createFileBrowser(name: string): HTMLDivElement {
  const input = createHiddenFileInput(name, onError);
  const body = createFileBrowserBody(onClick);

  let error = document.createElement("div");
  body.appendChild(error);

  const div = document.createElement("div");
  div.classList.add("file-browser");
  div.appendChild(body);
  div.appendChild(input);

  // Turn div into a drop target.
  // See https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/Drag_operations#specifying_drop_targets.
  div.addEventListener("dragstart", (event) => {
    event.preventDefault();
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = "copy";
    }
  });
  div.addEventListener("dragover", (event) => {
    event.preventDefault();
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = "copy";
    }
  });
  div.addEventListener("drop", (event) => {
    event.preventDefault();
    if (event.dataTransfer == null) {
      return;
    }

    // Upload file.
    const file = event.dataTransfer.files[0];
    uploadFile(name, file, onError);
  });
  return div;

  function onClick() {
    input.click();
  }

  function onError(message: string) {
    const replacement = createErrorMessage(message);
    error.replaceWith(replacement);
    error = replacement;
  }
}

// Wrapper around `uploadCSVFile` that checks for file validity and refreshes
// the page after a successful upload.
async function uploadFile(
  name: string,
  file: File | undefined,
  onError: (message: string) => void
) {
  if (file == null) {
    onError("Something went wrong. Please try again.");
    return;
  }
  if (file.type !== "text/csv") {
    onError("Not a CSV file.");
    return;
  }
  if (isTooBig(file)) {
    onError("The file is too big.");
    return;
  }
  await uploadCSVFile(name, file);

  // TODO refresh only if upload successful.
  window.location.href = window.location.href;
}
