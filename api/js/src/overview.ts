import { createActivityChart, createVocabularyChart } from "./chart";
import { getL1, getL2 } from "./language";
import { Activity } from "./schema";

function createOverviewHeader(): HTMLHeadingElement {
    const l1 = getL1();
    const l2 = getL2();
    const h1 = document.createElement("h1");
    h1.textContent = `${l2.name} from ${l1.name}`;
    return h1;
}

export function createOverviewPage(activityHistory: Activity[]): DocumentFragment {
    const fragment = document.createDocumentFragment();
    fragment.append(
        createOverviewHeader(),
        createVocabularyChart(activityHistory),
        createActivityChart(activityHistory),
    );
    return fragment;
}
