import {
    computeVocabularySize,
    createActivityChart,
    createVocabularyChart,
} from "./chart";
import { getL1, getL2 } from "./language";
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

    p.innerHTML = `
        <button is="button-link" href="/study">
            <img src="/public/svg/brain.svg?t=20221114"> Continue learning
        </button>
        <button is="button-link" href="/vocab">
            <img src="/public/svg/notebook.svg?t=20221114"> Vocabulary
        </button>
    `;
    return p;
}

function createTodaySummary(activityHistory: ActivityHistory): DocumentFragment {
    const { learned, strengthened, forgotten } = activityHistory.activities[0];
    const score = 100 * (learned + strengthened) / (learned + strengthened + forgotten);
    const template = document.createElement("template");
    template.innerHTML = `
        <h2>Recent activity</h2>
        <p>Summary of today's work:</p>
        <ul>
            <li>Learned ${learned} words</li>
            <li>Strengthened ${strengthened} words</li>
            <li>Forgot ${forgotten} words</li>
            <li>Your score: ${score}%</li>
        </ul>
    `;
    return template.content;
}

function hasActivity({ crammed, learned, strengthened }: Activity): boolean {
    return crammed > 0 || learned > 0 || strengthened > 0;
}

// Tries to compute streak.
// Since ActivityHistory only keeps track of activity in the past year,
// result may be less than the real streak.
function computeStreak(activityHistory: ActivityHistory): number {
    if (activityHistory.activities.length === 0) {
        return 0;
    }

    let streak = 0;
    for (let i = 1; i < activityHistory.activities.length; i++) {
        if (!hasActivity(activityHistory.activities[i])) {
            break;
        }
        streak = i;
    }
    if (hasActivity(activityHistory.activities[0])) {
        streak++;
    }
    return streak;
}

function createStreakSummary(activityHistory: ActivityHistory): DocumentFragment {
    const streak = computeStreak(activityHistory);
    const template = document.createElement("template");
    template.innerHTML = `
        <p>You're on a ${streak}-day streak.</p>
        <p class="button-group" style="justify-content: center">
            <button is="button-link" href="/study">
                <img src="/public/svg/heartbeat.svg?t=20221114"> Extend streak
            </button>
        </p>
    `;
    return template.content;
}

export function createOverviewPage(activityHistory: ActivityHistory): DocumentFragment {
    const fragment = document.createDocumentFragment();
    fragment.append(
        createOverviewHeader(),
        createVocabularyChart(activityHistory),
        createVocabularySummary(activityHistory),
        createActionButtons(),
        createTodaySummary(activityHistory),
        createActivityChart(activityHistory),
        createStreakSummary(activityHistory),
    );
    return fragment;
}
