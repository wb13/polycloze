import { VocabularyDataSchema } from "./schema";

import {
    BarController,
    BarElement,
    CategoryScale,
    Chart,
    LinearScale,
    PointElement,
} from "chart.js";

Chart.register(
    BarController,
    BarElement,
    CategoryScale,
    LinearScale,
    PointElement,
);

Chart.defaults.font.family = "Nunito";

const stackedBarChartOptions = {
    responsive: true,
    scales: {
        x: {
            stacked: true,
        },
        y: {
            stacked: true,
        },
    },
};

// NOTE supply your own data
const stackedBarChartConfig = {
    type: "bar",
    options: stackedBarChartOptions,
    plugins: {
        title: {
            display: true,
            text: "Vocabulary size",
        },
        legend: {
            labels: {
                font: {
                    family: "Nunito",
                },
            },
        },
    },
};

const monthLabels = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
];

function createChart(canvas: HTMLCanvasElement, vocabData: VocabularyDataSchema): Chart {
    const data = {
        labels: monthLabels,
        datasets: normalize(vocabData),
    };
    return new Chart(canvas, {
        ...stackedBarChartConfig,
        data,
    });
}

function randomColor(): string {
    return "#" + Math.floor(Math.random() * 0x10).toString(16).repeat(3);
}

// Makes data fit for chart.js use.
function normalize(data: VocabularyDataSchema): Array<{ label: string, data: number[] }> {
    const datasets = [];
    for (const dataset of data.datasets) {
        const label = dataset.name;
        const data = dataset.data.slice(0);
        datasets.push({
            label,
            data: label === "0h" ? data.map(x => -x) : data,
            backgroundColor: randomColor(),
        });
    }
    return datasets;
}

export function createVocabularyChart(data: VocabularyDataSchema): HTMLCanvasElement {
    const canvas = document.createElement("canvas");
    createChart(canvas, data);
    return canvas;
}
