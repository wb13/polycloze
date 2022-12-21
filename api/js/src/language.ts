// Language info stored in meta tags.

import { Language } from "./schema";

// Returns L1 metadata stored in meta tags.
export function getL1(): Language {
  // Guaranteed by the server to exist.
  const meta = document.querySelector(
    "meta[name='polycloze-l1']"
  ) as HTMLMetaElement;
  const { code, name, bcp47 } = meta.dataset;
  return { code: code as string, name: name as string, bcp47: bcp47 as string };
}

// Returns L2 metadata stored in meta tags.
export function getL2(): Language {
  // Guaranteed by the server to exist.
  const meta = document.querySelector(
    "meta[name='polycloze-l2']"
  ) as HTMLMetaElement;
  const { code, name, bcp47 } = meta.dataset;
  return { code: code as string, name: name as string, bcp47: bcp47 as string };
}
