import "./upload.css";
import { createButton } from "./button";
import { csrf } from "./csrf";
import { createLabeledIcon } from "./icon";
import { getL1, getL2 } from "./language";
import { resolve } from "./request";

// Checks if the file the user wants to upload is too big.
function isTooBig(file: File): boolean {
  // Returns true if file size > 8MB.
  return file.size > 8 * 1024 * 1024;
}

// Creates a hidden file input element for uploading CSV files.
function createHiddenFileInput(
  name: string,
  onSuccess: () => void,
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
    for (const file of files) {
      if (isTooBig(file)) {
        input.value = ""; // Clears the file list.
        onError("The file is too big.");
        return;
      }
    }
    onSuccess();
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

// Creates form data to submit to server.
// Includes a csrf token.
function createFormData(): FormData {
  const formData = new FormData();
  formData.append("csrf-token", csrf());
  return formData;
}

export function createFileBrowser(name: string): HTMLDivElement {
  const input = createHiddenFileInput(name, onSuccess, onError);
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
    if (file == null) {
      onError("Something went wrong. Please try again.");
      return;
    }
    const formData = createFormData();
    formData.append(name, file);

    const l1 = getL1().code;
    const l2 = getL2().code;

    const request = new XMLHttpRequest();
    request.open("POST", resolve(`/api/settings/upload/${l1}/${l2}`));
    request.send(formData);
    // TODO this works, but destroys success message, because redirect page
    // gets thrown away
  });
  return div;

  function onClick() {
    input.click();
  }

  function onSuccess() {
    // Submits files by clicking on submit button.
    // This gets triggers only when using the "Browse files" button.
    const button = createButton("");
    button.type = "submit";
    button.style.display = "none";
    div.appendChild(button);
    button.click();
  }

  function onError(message: string) {
    const replacement = createErrorMessage(message);
    error.replaceWith(replacement);
    error = replacement;
  }
}
