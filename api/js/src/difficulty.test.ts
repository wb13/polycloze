import { DifficultyTuner } from "./difficulty";

import assert from "assert";

describe("DifficultyTuner", () => {
  describe("constructor", () => {
    it("default correct and incorrect count should be 0", () => {
      const tuner = new DifficultyTuner();
      const { correct, incorrect } = tuner;

      assert.strict.equal(correct, 0);
      assert.strict.equal(incorrect, 0);
    });
  });

  describe("update", () => {
    it("update(true) should increase correct count or level", () => {
      const tuner = new DifficultyTuner();
      const { correct: correctBefore, level: levelBefore } = tuner;

      tuner.update(true);

      const { correct: correctAfter, level: levelAfter } = tuner;
      assert.ok(correctBefore < correctAfter || levelBefore < levelAfter);
    });

    it("update(false) should increase incorrect count or decrease level", () => {
      const tuner = new DifficultyTuner();
      const { incorrect: incorrectBefore, level: levelBefore } = tuner;

      tuner.update(false);

      const { incorrect: incorrectAfter, level: levelAfter } = tuner;
      assert.ok(incorrectBefore < incorrectAfter || levelAfter < levelBefore);
    });
  });
});
