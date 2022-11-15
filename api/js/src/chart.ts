import { fetchVocabularySize } from "./api";
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

const dayLabels = [
    "Mon",
    "Tue",
    "Wed",
    "Thu",
    "Fri",
    "Sat",
    "Sun",
];

function createChart(canvas: HTMLCanvasElement, vocabData: VocabularyDataSchema): Chart {
    const data = {
        labels: dayLabels,
        datasets: normalize(vocabData),
    };
    return new Chart(canvas, {
        type: "bar",
        options: stackedBarChartOptions,
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

type Timescale = "week";

export async function createVocabularyChart(_: Timescale): Promise<HTMLCanvasElement> {
    const vocab = await fetchVocabularySize({
        start: Math.floor(Date.now()/1000 - 7 * 24 * 3600),
        nSamples: 7,
    });

    const canvas = document.createElement("canvas");
    createChart(canvas, vocab);
    return canvas;
}
