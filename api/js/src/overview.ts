import {
    computeVocabularySize,
    createActivityChart,
    createVocabularyChart,
} from "./chart";
import { getL1, getL2 } from "./language";
import { createLink } from "./link";
import { Activity, ActivityHistory } from "./schema";

function createOverviewHeader(): HTMLHeadingElement {
    const l1 = getL1();
    const l2 = getL2();
    const h1 = document.createElement("h1");
    const title = `${l2.name} from ${l1.name}`;
    h1.textContent = title;
    document.title = `${title} | polycloze`;
    return h1;
}

function createVocabularySummary(activityHistory: ActivityHistory): HTMLParagraphElement {
    const size = computeVocabularySize(activityHistory)[0];
    const p = document.createElement("p");
    p.textContent = `You've learned ${size} words. Keep up the good work!`;
    return p;
}

function createActionButtons(): HTMLParagraphElement {
    const p = document.createElement("p");
    p.classList.add("button-group");
    p.style.justifyContent = "center";
    p.append(
        createLink("brain", "Continue learning", "/study"),
        createLink("notebook", "Vocabulary", "/vocab"),
    );
    return p;
}

function hasActivity({ crammed, learned, strengthened }: Activity): boolean {
    return crammed > 0 || learned > 0 || strengthened > 0;
}

// Tries to compute streak.
// Since ActivityHistory only keeps track of activity in the past year,
// result may be less than the real streak.
// Returns length of streak and boolean value (whether or not streak is active).
function computeStreak(activityHistory: ActivityHistory): [number, boolean] {
    if (activityHistory.activities.length === 0) {
        return [0, false];
    }

    let streak = 0;
    let active = false;
    for (let i = 1; i < activityHistory.activities.length; i++) {
        if (!hasActivity(activityHistory.activities[i])) {
            break;
        }
        streak = i;
    }
    if (hasActivity(activityHistory.activities[0])) {
        streak++;
        active = true;
    }
    return [streak, active];
}

function createStreakSummary(activityHistory: ActivityHistory): HTMLParagraphElement {
    const [streak, active] = computeStreak(activityHistory);

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

export function createOverviewPage(activityHistory: ActivityHistory): DocumentFragment {
    const h2 = document.createElement("h2");
    h2.textContent = "Recent activity";

    const fragment = document.createDocumentFragment();
    fragment.append(
        createOverviewHeader(),
        createVocabularyChart(activityHistory),
        createVocabularySummary(activityHistory),
        createActionButtons(),
        h2,
        createActivityChart(activityHistory),
        createStreakSummary(activityHistory),
    );
    return fragment;
}
