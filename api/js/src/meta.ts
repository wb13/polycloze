// Functions for accessing page metadata.

export type Metadata = {
  l1?: string;
  l2?: string;
};

export function getMetadata(): Metadata {
  const meta = document.querySelector(
    "meta[name='application-name'][value='polycloze']"
  ) as HTMLMetaElement | null;
  return {
    l1: meta?.dataset.l1Code,
    l2: meta?.dataset.l2Code,
  };
}
