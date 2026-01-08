import "./app.css";
import { ItemBuffer, RandomSentenceBuffer } from "./buffer";
import { createEmptyItem, createItem, createItemListen } from "./item";
import { TTS } from "./tts";

export async function createApp(
  buffer: ItemBuffer
): Promise<[HTMLDivElement, () => void]> {
  const tts = new TTS();
  await tts.init();

  const div = document.createElement("div");
  const item = await buffer.take();

  if (item == null) {
    return [createEmptyItem(), () => undefined];
  }

  const next = () => {
    createApp(buffer).then(([replacement, ready]) => {
      div.replaceWith(replacement);
      ready();
    });
  };

  const [child, resize] = createItem(tts, item, next);
  div.appendChild(child);

  const ready = () => {
    const blank = div.querySelector(".blank") as HTMLInputElement;
    blank.focus();
    resize();
  };
  return [div, ready];
}

export async function createListenApp(
  buffer: RandomSentenceBuffer
): Promise<[HTMLDivElement, () => void]> {
  const tts = new TTS();
  await tts.init();

  const div = document.createElement("div");
  const sentence = await buffer.take();

  if (sentence == null) {
    return [createEmptyItem(), () => undefined];
  }

  const parts = [{text: sentence.text, answers: [{text: sentence.text, normalized: sentence.text, new: false, difficulty: 0}]}];
  const item = {sentence: {id: sentence.id, tatoebaID: sentence.tatoebaID, parts: parts}, translation: sentence.translation};

  const next = () => {
    createListenApp(buffer).then(([replacement, ready]) => {
      div.replaceWith(replacement);
      ready();
    });
  };

  const [child, resize] = createItemListen(tts, item, next);
  div.appendChild(child);

  const ready = () => {
    const blank = div.querySelector(".blank") as HTMLInputElement;
    blank.focus();
    //resize();
  };
  return [div, ready];
}
