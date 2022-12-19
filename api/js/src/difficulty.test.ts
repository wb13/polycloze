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

      const changed = tuner.update(true);

      const { correct: correctAfter, level: levelAfter } = tuner;
      assert.ok(
        (!changed && correctBefore < correctAfter) ||
          (changed && levelBefore < levelAfter)
      );
    });

    it("update(false) should increase incorrect count or decrease level", () => {
      const tuner = new DifficultyTuner();
      const { incorrect: incorrectBefore, level: levelBefore } = tuner;

      const changed = tuner.update(false);

      const { incorrect: incorrectAfter, level: levelAfter } = tuner;
      assert.ok(
        (!changed && incorrectBefore < incorrectAfter) ||
          (changed && levelAfter < levelBefore)
      );
    });

    describe("level too easy", () => {
      describe("already at max level", () => {
        it("level shouldn't change", () => {
          const tuner = new DifficultyTuner({
            level: 5,
            max: 5,
            correct: 100, // too easy
          });

          const changed = tuner.update(true);
          assert.ok(!changed);
          assert.strict.equal(tuner.level, 5);
        });
      });

      describe("not yet at max level", () => {
        it("level should increase by 1", () => {
          const tuner = new DifficultyTuner({
            level: 4,
            max: 5,
            correct: 100, // too easy
          });

          const changed = tuner.update(true);
          assert.ok(changed);
          assert.strict.equal(tuner.level, 5);
        });
      });
    });

    describe("level too hard", () => {
      describe("already at min level", () => {
        it("level shouldn't change", () => {
          const tuner = new DifficultyTuner({
            level: 1,
            min: 1,
            incorrect: 100, // too hard
          });

          const changed = tuner.update(false);
          assert.ok(!changed);
          assert.strict.equal(tuner.level, 1);
        });
      });

      describe("not yet at min level", () => {
        it("level should decrease by 1", () => {
          const tuner = new DifficultyTuner({
            level: 1,
            min: 0,
            incorrect: 100, // too hard
          });

          const changed = tuner.update(false);
          assert.ok(changed);
          assert.strict.equal(tuner.level, 0);
        });
      });
    });
  });
});
