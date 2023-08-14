import { ActivitySummary, DataPoint } from "./schema";

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
  // @ts-ignore
  Colors,
  Filler,
  Legend,
  LineController,
  LineElement,
  LinearScale,
  PointElement,
  Tooltip
);

Chart.defaults.font.family = "Nunito";

const dayLabels = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

// Returns vocabulary size data over the past week.
// Assumes that `vocabularySize` has a data point for each day in the past
// week.
function vocabularyData(
  vocabularySize: DataPoint[],
  estimatedLevel: DataPoint[]
): ChartData {
  const points = vocabularySize.slice(-7);
  const dataVocab = points.map((point) => point.value);
  const dataEstimate = estimatedLevel
    .slice(-7)
    .map((point) => Math.floor(0.85 * 2 ** (point.value - 1)));
  // Difficulty tuner increases difficulty when it's 85% confident that the
  // student already knows 85% of all words in the current difficulty.
  const labels = points.map((point) => dayLabels[point.time.getDay()]);
  return {
    labels,
    datasets: [
      {
        data: dataVocab,
        label: "Observed vocabulary size",
        cubicInterpolationMode: "monotone",
        fill: true,
      },
      {
        data: dataEstimate,
        label: "Estimated vocabulary size",
        cubicInterpolationMode: "monotone",
        fill: true,
      },
    ],
  };
}

// Creates dataset for activity data.
function createDataset(label: string, data: number[]): ChartDataset {
  return { data, label, cubicInterpolationMode: "monotone" };
}

// Formats activity data so it can be used with chart.js.
// Assumes `activity` has a record for each day in the past week.
function activityData(activity: ActivitySummary[]): ChartData {
  const points = activity.slice(-7);
  const labels = points.map((point) => dayLabels[point.from.getDay()]);
  return {
    labels,
    datasets: [
      createDataset(
        "Learned",
        points.map((p) => p.learned)
      ),
      createDataset(
        "Forgotten",
        points.map((p) => p.forgotten)
      ),
      createDataset(
        "Unimproved",
        points.map((p) => p.unimproved)
      ),
      createDataset(
        "Crammed",
        points.map((p) => p.crammed)
      ),
      createDataset(
        "Strengthened",
        points.map((p) => p.strengthened)
      ),
    ],
  };
}

function createChart(
  canvas: HTMLCanvasElement,
  vocabularySize: DataPoint[],
  estimatedLevel: DataPoint[]
): Chart {
  return new Chart(canvas, {
    type: "line",
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: {
          grid: {
            color: "#444"
          }
        },
        y: {
          min: 0,
          grid: {
            color: "#444"
          }
        },
      },
      plugins: {
        legend: {
          position: "bottom",
        },
      },
    },
    data: vocabularyData(vocabularySize, estimatedLevel),
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
  estimatedLevel: DataPoint[]
): HTMLDivElement {
  const canvas = document.createElement("canvas");
  createChart(canvas, vocabularySize, estimatedLevel);
  return createChartContainer(canvas);
}

export function createActivityChart(
  activity: ActivitySummary[]
): HTMLDivElement {
  const canvas = document.createElement("canvas");
  new Chart(canvas, {
    type: "line",
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: {
          grid: {
            color: "#444"
          }
        },
        y: {
          min: 0,
          grid: {
            color: "#444"
          }
        },
      },
      plugins: {
        legend: {
          position: "bottom",
        },
      },
    },
    data: activityData(activity),
  });
  return createChartContainer(canvas);
}
