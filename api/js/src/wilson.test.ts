import { isTooEasy, isTooHard } from "./wilson";

import assert from "assert";

describe("isTooEasy", () => {
  describe("0 correct, 0 incorrect", () => {
    it("should return false", () => {
      assert.ok(!isTooEasy(0, 0));
    });
  });

  describe("lots of correct answers", () => {
    it("should return true", () => {
      assert.ok(isTooEasy(100, 1));
    });
  });
});

describe("isTooHard", () => {
  describe("0 correct, 0 incorrect", () => {
    it("should return false", () => {
      assert.ok(!isTooHard(0, 0));
    });
  });

  describe("lots of incorrect answers", () => {
    it("should return true", () => {
      assert.ok(isTooHard(0, 4));
    });
  });
});
