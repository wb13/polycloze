import { substituteDigraphs } from "./digraph";

import assert from "assert";

describe("substituteDigraphs", () => {
  describe("text has no digraphs", () => {
    it("should return input unchanged", () => {
      const text = "hello, world!";
      assert.strict.equal(text, substituteDigraphs(text));
    });
  });

  describe("text has digraphs", () => {
    it("should transform all digraphs in the text", () => {
      const text = "Meine \\A:pfel sind auch deine \\A:pfel.";
      const result = "Meine Äpfel sind auch deine Äpfel.";
      assert.strict.equal(substituteDigraphs(text), result);
    });
  });

  describe("text has backslashes but aren't followed by a valid digraph", () => {
    it("should return input unchanged", () => {
      const text = "\\hhh\\";
      assert.strict.equal(text, substituteDigraphs(text));
    });
  });
});
