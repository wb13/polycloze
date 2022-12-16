// Prompts user if they really want to leave the page if there are unsaved changes.

let unsaved = 0;

// Usage:
//
// ```
// const save = edit();
// // do something (i.e. notify server)
// save();
// ```
export function edit(): () => void {
  let done = false;
  ++unsaved;
  return () => {
    if (!done) {
      --unsaved;
      done = true;
    }
  };
}

addEventListener("beforeunload", (event) => {
  if (unsaved > 0) {
    event.preventDefault();
    return (event.returnValue = "");
  }
});
