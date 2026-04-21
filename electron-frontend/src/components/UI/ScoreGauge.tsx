import React from 'react';

interface ScoreGaugeProps {
  score: number;
  label?: string;
}

export const ScoreGauge: React.FC<ScoreGaugeProps> = ({ score, label }) => {
  return (
    <div className="score-gauge">
      {label && <span>{label}: </span>}
      <span>{score}%</span>
    </div>
  );
};