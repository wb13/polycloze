self.addEventListener("install", (event: Event) => {
  console.log("install:", event);
});

self.addEventListener("activate", (event: Event) => {
  console.log("activate:", event);
});
