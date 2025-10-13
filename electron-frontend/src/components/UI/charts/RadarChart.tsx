import '../../../Styles/Components/UI/RadarChart.css';

import {
  Chart as ChartJS,
  Filler,
  Legend,
  LineElement,
  PointElement,
  RadialLinearScale,
  Tooltip,
} from 'chart.js';

import { Radar } from 'react-chartjs-2';
import React from 'react';

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
);

interface Metric {
  score_title: string;
  weighted_score: string;
  raw_score: string; 
  score_reason: string;
}

interface RadarChartOverviewProps {
  metrics: Metric[];
}

const RadarChartOverview: React.FC<RadarChartOverviewProps> = ({ metrics }) => {
  const weightedData = metrics.map((m) => Number(m.weighted_score));
  const unweightedData = metrics.map((m) => Number(m.raw_score));
  console.log(weightedData);
  console.log(unweightedData);

  const data = {
    labels: metrics.map((m) => m.score_title),
    datasets: [
      {
        label: 'Weighted Match Score',
        data: weightedData,
        backgroundColor: 'rgba(46, 204, 113, 0.3)', 
        borderColor: 'rgba(43, 122, 120, 0.8)', 
        borderWidth: 2,
        pointBackgroundColor: 'rgba(43, 122, 120, 0.8)',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: 'rgba(43, 122, 120, 1)',
        fill: true,
      },
      {
        label: 'Raw Match Score',
        data: unweightedData,
        backgroundColor: 'rgba(52, 152, 219, 0.3)', 
        borderColor: 'rgba(43, 122, 120, 0.8)',
        borderWidth: 2,
        pointBackgroundColor: 'rgba(43, 122, 120, 0.8)',
        pointBorderColor: '#fff',
        pointHoverBackgroundColor: '#fff',
        pointHoverBorderColor: 'rgba(43, 122, 120, 1)',
        fill: true,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false, 
    scales: {
      r: {
        min: 0,
        max: 100,
        angleLines: {
          color: '#555', 
        },
        grid: {
          color: '#555', 
          circular: true, 
        },
        pointLabels: {
          color: '#fff', 
          font: { size: 14 },
        },
        ticks: {
          display: false, 
          beginAtZero: true,
          stepSize: 5,
        },
      },
    },
    plugins: {
      legend: { labels: { color: '#fff' } },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            const idx = context.dataIndex;
            return `${metrics[idx].score_reason} (${context.formattedValue}%)`;
          },
        },
      },
    },
  };
  /*  */

  return <Radar data={data} options={options} />;
};

export default RadarChartOverview;
