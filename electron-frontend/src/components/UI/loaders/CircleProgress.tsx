import '@/assets/styles/Components/UI/loaders/CircleProgress.css';

import React, { useEffect, useState } from 'react';

interface CircleProgressProps {
  percentage?: number;
  size?: number;
  strokeWidth?: number;
  tickCount?: number;
}

export const CircleProgress: React.FC<CircleProgressProps> = ({
  percentage,
  size = 200,
  strokeWidth = 6,
  tickCount = 100,
}) => {
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    if (percentage === undefined) {
      setProgress(0);
      return;
    }

    let start = 0;
    const interval = setInterval(() => {
      const step = Math.max(0.5, (percentage - start) / 10);
      start += step;

      if (start >= percentage) {
        start = percentage;
        clearInterval(interval);
      }
      setProgress(Math.round(start));
    }, 20);

    return () => clearInterval(interval);
  }, [percentage]);

  const radius = size / 2 - strokeWidth;
  const center = size / 2;

  if (percentage === undefined) {
    return (
      <div style={{ position: 'relative', width: size, height: size }}>
        <svg className="circle-progress-indeterminate" viewBox="25 25 50 50">
          <circle cx="50" cy="50" r="20" fill="none" strokeWidth={strokeWidth}></circle>
        </svg>
      </div>
    );
  }

  const getTickColor = (i: number) => {
    const ratio = i / tickCount;
    const red = Math.round(255 * (1 - ratio));
    const green = Math.round(255 * ratio);
    return `rgb(${red}, ${green}, 0)`;
  };

  const ticks = Array.from({ length: tickCount }, (_, i) => {
    const angle = (i / tickCount) * 360;
    const isActive = i < progress;
    const transform = `rotate(${angle}deg) translate(${radius}px)`;
    return (
      <div
        key={i}
        style={{
          position: 'absolute',
          top: center - 1,
          left: center - 1,
          width: 2,
          height: 8,
          backgroundColor: isActive ? getTickColor(i) : '#444',
          transformOrigin: '0 0',
          transform,
        }}
      />
    );
  });

  return (
    <div style={{ position: 'relative', width: size, height: size }}>
      {ticks}
      <div
        style={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: size,
          height: size,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          flexDirection: 'column',
          fontSize: size * 0.25,
          fontWeight: 'bold',
          color: '#fff',
        }}
      >
        <div style={{ fontSize: size * 0.12, fontWeight: 'normal' }}>
          Generating...
        </div>
        <div style={{ fontSize: size * 0.25, fontWeight: 'bold' }}>
          {progress}%
        </div>
      </div>
    </div>
  );
};