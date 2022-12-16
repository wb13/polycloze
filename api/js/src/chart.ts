import { ActivityHistory, DataPoint } from "./schema";

import {
    CategoryScale,
    Chart,
    ChartData,
    ChartDataset,
    Colors,
    Filler,
    Legend,
    LineController,
    LineElement,
    LinearScale,
    PointElement,
    Tooltip,
} from "chart.js";

Chart.register(
    CategoryScale,
    Colors,
    Filler,
    Legend,
    LineController,
    LineElement,
    LinearScale,
    PointElement,
    Tooltip,
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
// Assumes that `vocabularySize` has a data point for each day in the past
// week.
function vocabularyData(vocabularySize: DataPoint[]): ChartData {
    const points = vocabularySize.slice(-7);
    const data = points.map(point => point.value);
    const labels = points.map(point => dayLabels[point.time.getDay()]);
    return {
        labels,
        datasets: [
            {
                data,
                label: "Vocabulary size",
                cubicInterpolationMode: "monotone",
                fill: true,
            },
        ],
    };
}

// Creates dataset for activity data.
function createDataset(label: string, data: number[]): ChartDataset {
    return {data, label, cubicInterpolationMode: "monotone"};
}

function activityData(activityHistory: ActivityHistory): ChartData {
    const week = activityHistory.activities.slice(0, 7);
    const labels = week.map((_, i) => dayLabels[dateNDaysAgo(i).getDay()]).reverse();
    return {
        labels,
        datasets: [
            createDataset("Learned", week.map(a => a.learned).reverse()),
            createDataset("Forgotten", week.map(a => a.forgotten).reverse()),
            createDataset("Unimproved", week.map(a => a.unimproved).reverse()),
            createDataset("Crammed", week.map(a => a.crammed).reverse()),
            createDataset("Strengthened", week.map(a => a.strengthened).reverse()),
        ],
    };
}

function createChart(
    canvas: HTMLCanvasElement,
    vocabularySize: DataPoint[],
): Chart {
    return new Chart(canvas, {
        type: "line",
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    min: 0,
                    ticks: {
                        stepSize: 1,
                    },
                },
            },
            plugins: {
                legend: {
                    position: "bottom",
                },
            },
        },
        data: vocabularyData(vocabularySize),
    });
}

// Wraps around chart to make it responsive.
// See https://www.chartjs.org/docs/latest/configuration/responsive.html.
function createChartContainer(chart: HTMLCanvasElement): HTMLDivElement {
    const div = document.createElement("div");
    div.style.position = "relative";
    div.appendChild(chart);
    return div;
}

export function createVocabularyChart(
    vocabularySize: DataPoint[],
): HTMLDivElement {
    const canvas = document.createElement("canvas");
    createChart(canvas, vocabularySize);
    return createChartContainer(canvas);
}

export function createActivityChart(activityHistory: ActivityHistory): HTMLDivElement {
    const canvas = document.createElement("canvas");
    new Chart(canvas, {
        type: "line",
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: "bottom",
                },
            },
            scales: {
                y: {
                    ticks: {
                        stepSize: 1,
                    },
                },
            },
        },
        data: activityData(activityHistory),
    });
    return createChartContainer(canvas);
}
