import { createActivityChart, createVocabularyChart } from "./chart";
import { startOfDay, endOfDay } from "./datetime";
import { getL1, getL2 } from "./language";
import { createLink } from "./link";
import { Activity, ActivitySummary, DataPoint } from "./schema";

function createOverviewHeader(): HTMLHeadingElement {
  const l1 = getL1();
  const l2 = getL2();
  const h1 = document.createElement("h1");
  const title = `${l2.name} from ${l1.name}`;
  h1.textContent = title;
  document.title = `${title} | polycloze`;
  return h1;
}

// Gives words of encouragement to the user.
function encourage(): string {
  const phrases = ["Why don't you give it a try?", "You can do it!"];
  const choice = Math.floor(Math.random() * phrases.length);
  return phrases[choice];
}

// Praises the user.
function praise(): string {
  const phrases = ["Keep up the good work!", "Keep it up!", "Good work!"];
  const choice = Math.floor(Math.random() * phrases.length);
  return phrases[choice];
}

function createVocabularySummary(vocabularySize: number): HTMLParagraphElement {
  const p = document.createElement("p");
  if (vocabularySize <= 0) {
    p.textContent = `You haven't learned any words yet. ${encourage()}`;
  } else {
    p.textContent = `You've learned ${vocabularySize} words. ${praise()}`;
  }
  return p;
}

function createActionButtons(vocabularySize: number): HTMLParagraphElement {
  const text = vocabularySize <= 0 ? "Start learning" : "Continue learning";

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.style.justifyContent = "center";
  p.append(
    createLink("brain", text, "/study"),
    createLink("notebook", "Vocabulary", "/vocab")
  );
  return p;
}

function hasActivity({ crammed, learned, strengthened }: Activity): boolean {
  return crammed > 0 || learned > 0 || strengthened > 0;
}

// Checks if the date intervals overlap.
// Intervals are assumed to be half-open: [start, end).
function isOverlapping(a1: Date, a2: Date, b1: Date, b2: Date): boolean {
  // Rearrange args so that x1 < x2.
  if (a2 < a1) {
    [a1, a2] = [a2, a1];
  }
  if (b2 < b1) {
    [b1, b2] = [b2, b1];
  }

  // Check for overlaps.
  if (a2 <= b1) {
    return false;
  }
  if (b2 <= a1) {
    return false;
  }
  return true;
}

// Tries to compute current streak and if it's active or not.
// The returned value may be shorter than the actual streak, if the streak
// started before the earliest day in `activity`.
function computeActiveStreak(activity: ActivitySummary[]): [number, boolean] {
  let active = false;

  // Construct range for today.
  const now = new Date();
  const todayStart = startOfDay(now);
  const todayEnd = endOfDay(now);

  // Find start of streak by looking for last day without activity,
  // excluding today.
  let start = now;
  for (let i = activity.length - 1; i >= 0; i--) {
    const summary = activity[i];
    const isToday = isOverlapping(
      todayStart,
      todayEnd,
      summary.from,
      summary.to
    );
    if (isToday) {
      if (hasActivity(summary)) {
        active = true;
        start = summary.from;
      }
      continue;
    }
    if (!hasActivity(summary)) {
      break;
    }
    start = summary.from;
  }

  const diff = now.valueOf() - start.valueOf();
  let streak = Math.max(0, Math.floor(diff / 1000 / 60 / 60 / 24));
  if (active) {
    streak++;
  }
  return [streak, active];
}

function createStreakSummary(
  streak: number,
  active: boolean
): HTMLParagraphElement {
  let icon: string;
  let message: string;
  if (streak <= 0) {
    icon = "barbell";
    message = "Practice today";
  } else if (active) {
    icon = "barbell";
    message = "Keep practicing";
  } else {
    icon = "heartbeat";
    message = `Extend your ${streak}-day streak`;
  }

  const p = document.createElement("p");
  p.classList.add("button-group");
  p.style.justifyContent = "center";
  p.append(createLink(icon, message, "/study"));
  return p;
}

export function createOverviewPage(
  activity: ActivitySummary[],
  vocabularySize: DataPoint[],
  estimatedLevel: DataPoint[]
): DocumentFragment {
  const size = vocabularySize[vocabularySize.length - 1].value;
  const [streak, active] = computeActiveStreak(activity);

  const h2 = document.createElement("h2");
  h2.textContent = "Recent activity";

  const fragment = document.createDocumentFragment();
  fragment.append(
    createOverviewHeader(),
    createVocabularyChart(vocabularySize, estimatedLevel),
    createVocabularySummary(size),
    createActionButtons(size),
    h2,
    createActivityChart(activity),
    createStreakSummary(streak, active)
  );
  return fragment;
}
