import { Activity } from "./schema";

import {
    CategoryScale,
    Chart,
    ChartData,
    ChartDataset,
    Filler,
    LineController,
    LineElement,
    LinearScale,
    PointElement,
} from "chart.js";

Chart.register(
    CategoryScale,
    Filler,
    LineController,
    LineElement,
    LinearScale,
    PointElement,
);

Chart.defaults.font.family = "Nunito";

const dayLabels = [
    "Sun",
    "Mon",
    "Tue",
    "Wed",
    "Thu",
    "Fri",
    "Sat",
];

// Returns date n days ago.
function dateNDaysAgo(n: number): Date {
    const daysSinceEpoch = Date.now() / 1000 / 60 / 60 / 24 - n;
    return new Date(daysSinceEpoch * 24 * 60 * 60 * 1000);
}

// Returns vocabulary size data over the past week.
function vocabularyData(activityHistory: Activity[]): ChartData {
    const week = activityHistory.slice(0, 7);
    const labels = week.map((_, i) => dayLabels[dateNDaysAgo(i).getDay()]).reverse();
    const data = week.map(a => a.learned - a.forgotten).reverse();
    for (let i = 1; i < data.length; i++) {
        data[i] += data[i-1];
    }
    return {
        labels,
        datasets: [{data, cubicInterpolationMode: "monotone", fill: true}],
    };
}

function createDataset(data: number[]): ChartDataset {
    return {data, cubicInterpolationMode: "monotone"};
}

function activityData(activityHistory: Activity[]): ChartData {
    const week = activityHistory.slice(0, 7);
    const labels = week.map((_, i) => dayLabels[dateNDaysAgo(i).getDay()]).reverse();
    return {
        labels,
        datasets: [
            createDataset(week.map(a => a.forgotten).reverse()),
            createDataset(week.map(a => a.unimproved).reverse()),
            createDataset(week.map(a => a.crammed).reverse()),
            createDataset(week.map(a => a.learned).reverse()),
            createDataset(week.map(a => a.strengthened).reverse()),
        ],
    };
}

function createChart(canvas: HTMLCanvasElement, activityHistory: Activity[]): Chart {
    return new Chart(canvas, {
        type: "line",
        options: {
            responsive: true,
        },
        data: vocabularyData(activityHistory),
    });
}

export function createVocabularyChart(activityHistory: Activity[]): HTMLCanvasElement {
    const canvas = document.createElement("canvas");
    createChart(canvas, activityHistory);
    return canvas;
}

export function createActivityChart(activityHistory: Activity[]): HTMLCanvasElement {
    const canvas = document.createElement("canvas");
    new Chart(canvas, {
        type: "line",
        options: {
            responsive: true,
        },
        data: activityData(activityHistory),
    });
    return canvas;
}
