import React, { useEffect, useMemo, useState } from 'react';
import { useSetHeaderControls, useSetHeaderSubtitle, useSetHeaderTitle } from "@/components/Layouts/providers/HeaderProvider";

import { ApplicationStatus } from '../types';
import { useApplication } from '../providers/ApplicationProvider';

interface ProcessedApplicantData {
  RoleID: number;
  ApplicationStatus: ApplicationStatus;
  matchScore: number | undefined; 
  hasSummary: boolean;
  hasCoverLetter: boolean;
  isVeteran: boolean;
  isDisabled: boolean;
  raceMentioned: boolean;
  isEasyApply: boolean;
  timeSubmitted: Date;
  interviewed: boolean;
  codeAssessmentRequested: boolean;
  jitterY: number;
}

const factorLabels: Partial<Record<keyof ProcessedApplicantData, string>> = {
  hasCoverLetter: "Included Cover Letter",
  isEasyApply: "Easy Apply Used",
  isVeteran: "Veteran Status",
  // isDisabled: "Disability Status",
  // raceMentioned: "Race Mentioned",
  hasSummary: "Included Summary",
};

export const ApplicantDataPage: React.FC = () => {
  const { jobs } = useApplication();
  
  const setHeaderTitle = useSetHeaderTitle();
  const setHeaderSubtitle = useSetHeaderSubtitle();
  const setHeaderControls = useSetHeaderControls();
  
  useEffect(() => {
    setHeaderTitle("Applicant Data Analytics");
    setHeaderSubtitle("Overview of your application performance metrics");
    setHeaderControls(null); // Clear any controls left over from other pages
  }, [setHeaderTitle, setHeaderSubtitle, setHeaderControls]);
  
  const [data, setData] = useState<ProcessedApplicantData[]>([]);
  
  useEffect(() => {
    const loadData = async () => {
      const processed = await Promise.all(jobs.map(async (job) => {
        const savedMetrics = localStorage.getItem(`metrics-${job.RoleID}`);
        const metrics = savedMetrics ? JSON.parse(savedMetrics) : {};
        
        const fileName = `match-summary/${job.CompanyProperName.toLowerCase().replaceAll(
          " ",
          "_"
        )}_${job.JobTitle.toLowerCase().replaceAll(" ", "_")}_${job.RoleID}_match_summary.json`;
        
        let score: number | undefined = undefined;
        try {
          if (window.appAPI?.files) {
            const result = await window.appAPI.files.readJsonFile(fileName);
            // If the file exists and was read successfully, extract the overall match score
            if (result.success && result.data?.match_summary?.overall_match_summary?.overall_match_score) {
              score = result.data.match_summary.overall_match_summary.overall_match_score;
            }
          }
        } catch (e) {
          console.error(`Failed to load match score for RoleID ${job.RoleID}`, e);
        }
        
        const stableJitter = ((job.RoleID * 43) % 10) - 5;
        
        return {
          RoleID: job.RoleID,
          ApplicationStatus: job.ApplicationStatus,
          matchScore: score,
          hasSummary: !!metrics.hasSummary,
          hasCoverLetter: !!metrics.hasCoverLetter,
          isVeteran: !!metrics.isVeteran,
          isDisabled: !!metrics.isDisabled,
          raceMentioned: !!metrics.raceMentioned,
          isEasyApply: !!metrics.isEasyApply,
          timeSubmitted: metrics.timeSubmitted ? new Date(metrics.timeSubmitted) : new Date(),
          interviewed: !!metrics.interviewed,
          codeAssessmentRequested: !!metrics.codeAssessmentRequested,
          jitterY: stableJitter,
        };
      }));
      
      setData(processed);
    };
    
    if (jobs.length > 0) {
      loadData();
    }
  }, [jobs]);
  
  const calculateWinRate = (factor: keyof ProcessedApplicantData) => {
    const groupA = data.filter(d => d[factor] === true);
    const groupB = data.filter(d => d[factor] === false);
    
    const getRate = (list: ProcessedApplicantData[]) => {
      const wins = list.filter(d => d.interviewed || d.codeAssessmentRequested).length;
      
      const losses = list.filter(d => 
        d.ApplicationStatus === "Rejected" && 
        !d.interviewed && 
        !d.codeAssessmentRequested
      ).length;
      
      const totalDecided = wins + losses;
      return totalDecided > 0 ? (wins / totalDecided) * 100 : 0;
    };
    
    const groupARate = getRate(groupA);
    const groupBRate = getRate(groupB);
    
    return {
      groupARate: groupARate.toFixed(1),
      groupBRate: groupBRate.toFixed(1),
      delta: (groupARate - groupBRate).toFixed(1),
    };
  };
  
  const factors: (keyof ProcessedApplicantData)[] = ['hasSummary', 'hasCoverLetter', 'isEasyApply'];
  
  const dayOfWeekData = useMemo(() => {
    const counts = new Array(7).fill(0);
    const today = new Date().setHours(0, 0, 0, 0);
    
    data.filter(d => d.interviewed || d.codeAssessmentRequested || d.ApplicationStatus === "Rejected").forEach(d => {
      const submittedDay = new Date(d.timeSubmitted).setHours(0, 0, 0, 0);
      if (submittedDay === today) return;
      
      const day = d.timeSubmitted.getDay();
      if (day >= 0 && day < 7) {
        counts[day]++;
      }
    });
    return counts;
  }, [data]);
  
  const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  
  return (
    <div style={{ padding: '20px', color: 'var(--text-main)', overflowY: 'auto', height: '100%' }}>
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '30px' }}>
    
    {/* Graph 1: Win-Rate Comparison */}
    <div style={{ background: 'var(--bg-secondary)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
    <h3>Advancement Rate Comparison </h3>
    <p style={{ fontSize: '0.9em', color: 'var(--text-muted)', marginBottom: '15px' }}>
    Impact of factors on Advance Rate (Positive Delta = Higher chance of Application Advancement).
    Delta represents the percentage point difference in advancement rate between Group A (True) and Group B (False) for each factor.
    </p>
    <table style={{ width: '100%', textAlign: 'left', borderCollapse: 'collapse' }}>
    <thead>
    <tr style={{ borderBottom: '1px solid var(--border-color)' }}>
    <th style={{ padding: '8px 0' }}>Factor</th>
    <th>Group A (True)</th>
    <th>Group B (False)</th>
    <th style={{ paddingLeft: '10px' }}>Delta</th>
    </tr>
    </thead>
    <tbody>
    {factors.map(factor => {
      const rates = calculateWinRate(factor);
      const deltaNum = Number.parseFloat(rates.delta);
      const isPositive = deltaNum > 0;
      
      return (
        <tr key={factor} style={{ borderBottom: '1px solid var(--border-color)' }}>
        <td style={{ padding: '8px 0' }}>{factorLabels[factor] || factor}</td>
        <td>{rates.groupARate}%</td>
        <td>{rates.groupBRate}%</td>
        <td style={{ paddingLeft: '10px' }}>
        <div style={{ display: 'flex', alignItems: 'center', minWidth: '160px' }}>
        <span style={{ 
          width: '55px',
          color: isPositive ? 'lightgreen' : 'salmon',
          fontWeight: 'bold',
          whiteSpace: 'nowrap',
          marginRight: '8px'
        }}>
        {isPositive ? '+' : ''}{rates.delta}%
        </span>
        
        <div style={{ 
          flex: 1,
          height: '10px', 
          background: 'var(--bg-main)', 
          position: 'relative'
        }}>
        <div style={{ 
          position: 'absolute', 
          height: '100%', 
          background: isPositive ? 'green' : 'red',
          width: `${Math.min(Math.abs(deltaNum) * 2, 50)}%`,
          [isPositive ? 'left' : 'right']: '50%',
        }} />
        </div>
        </div>
        </td>
        </tr>
      );
    })}
    </tbody>
    </table>
    </div>
    
    {/* Graph 2: Match Score vs Outcome Scatter Plot */}
    <div style={{ background: 'var(--bg-secondary)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
    <h3>Predictive Performance: Overall Match vs. Progression</h3>
    <p style={{ fontSize: '0.9em', color: 'var(--text-muted)' }}>
    A breakdown of advancement rates (Interviews/Assessments) relative to the internal Match Score. Scores are weighted by a combination of NLP-based keyword matching, applicant pool volume, and alignment with mandatory vs. preferred requirements.
    </p>
    <div style={{ position: 'relative', height: '200px', width: 'calc(100% - 85px)', marginLeft: '85px', borderLeft: '2px solid var(--text-main)', borderBottom: '2px solid var(--text-main)', marginTop: '20px', marginBottom: '40px' }}>
    {[0, 20, 40, 60, 80, 100].map(tick => (
      <div key={tick} style={{ position: 'absolute', left: `${tick}%`, bottom: '-5px', height: '5px', borderLeft: '2px solid var(--text-main)' }}>
      <span style={{ position: 'absolute', top: '8px', left: '50%', transform: 'translateX(-50%)', fontSize: '0.8em', color: 'var(--text-muted)' }}>
      {tick}
      </span>
      </div>
    ))}
    <div style={{ position: 'absolute', bottom: '-40px', left: '50%', transform: 'translateX(-50%)', fontSize: '0.85em', color: 'var(--text-main)', fontWeight: 'bold' }}>
    Match Score
    </div>
    {data.filter(d => 
      (d.matchScore !== undefined && d.matchScore !== null) && 
      ((d.interviewed || d.codeAssessmentRequested) || d.ApplicationStatus !== "Open")
    ).map((d) => {
      const advanced = d.interviewed || d.codeAssessmentRequested;
      return (
        <div
        key={d.RoleID}
        style={{
          position: 'absolute',
          left: `${d.matchScore}%`,
          bottom: advanced ? '70%' : '20%', 
          transform: `translate(-50%, ${d.jitterY}px)`, 
          width: '8px',
          height: '8px',
          borderRadius: '50%',
          background: advanced ? 'lightgreen' : 'salmon',
          opacity: 0.7
        }}
        title={`Score: ${d.matchScore}, Advanced: ${advanced}`}
        />
      );
    })}
    <span style={{ position: 'absolute', bottom: '70%', right: 'calc(100% + 10px)', fontSize: '0.9em' }}>Advanced</span>
    <span style={{ position: 'absolute', bottom: '20%', right: 'calc(100% + 10px)', fontSize: '0.9em' }}>Rejected</span>
    </div>
    </div>
    
    {/* Graph 3: Day of Week Analysis */}
    <div style={{ background: 'var(--bg-secondary)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
    <h3>Response Distribution by Weekday (Response Rate)</h3>
    <p style={{ fontSize: '0.9em', color: 'var(--text-muted)', marginBottom: '10px' }}>
    Identifies which application submission days yield the highest engagement. This data tracks interviews, assessments, or rejections to help target optimal windows for submission.
    </p>
    <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '20px' }}>
    {dayOfWeekData.map((count, dayIndex) => {
      const maxCount = Math.max(...dayOfWeekData, 1); 
      const widthPercent = (count / maxCount) * 100;
      const dayName = daysOfWeek[dayIndex];
      
      return (
        <div key={dayName} style={{ display: 'flex', alignItems: 'center' }}>
        <div style={{ width: '40px', fontSize: '0.8em' }}>{daysOfWeek[dayIndex]}</div>
        <div style={{ flex: 1, background: 'var(--bg-main)', height: '20px', borderRadius: '4px', overflow: 'hidden' }}>
        <div style={{ 
          width: `${widthPercent}%`, 
          height: '100%', 
          background: count > 0 ? `rgba(0, 255, 0, ${0.4 + (count/maxCount)*0.6})` : 'transparent',
          transition: 'width 0.3s ease'
        }} title={`${daysOfWeek[dayIndex]} - Responses: ${count}`} />
        </div>
        <div style={{ width: '30px', textAlign: 'right', fontSize: '0.8em' }}>{count}</div>
        </div>
      );
    })}
    </div>
    </div>
    
    {/* Graph 4: Logistic Regression Proxy */}
    <div style={{ background: 'var(--bg-secondary)', padding: '20px', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
    <h3>Predictive Power of Application Factors (95% Confidence Intervals)</h3>
    <p style={{ fontSize: '0.9em', color: 'var(--text-muted)' }}>
    This chart plots the estimated impact of each factor on the advancement rate. Narrow lines indicate high-confidence patterns, while wider lines suggest that more data is required. If a factor crosses the center line, its impact is currently indistinguishable from random chance.
    </p>
    <div style={{ position: 'relative', height: '230px', width: '100%', marginTop: '20px' }}>
    <div style={{ position: 'absolute', left: '50%', top: 0, height: '180px', borderLeft: '2px dashed var(--text-muted)', zIndex: 0 }} />
    <div style={{ position: 'absolute', top: '180px', left: 0, width: '100%', borderBottom: '2px solid var(--text-main)', zIndex: 0 }} />
    <span style={{ position: 'absolute', top: '190px', left: '0', fontSize: '0.75em', color: 'var(--text-muted)' }}>-50% Impact</span>
    <span style={{ position: 'absolute', top: '190px', left: '50%', transform: 'translateX(-50%)', fontSize: '0.75em', color: 'var(--text-main)', fontWeight: 'bold' }}>0 (Neutral)</span>
    <span style={{ position: 'absolute', top: '190px', right: '0', fontSize: '0.75em', color: 'var(--text-muted)' }}>+50% Impact</span>
    
    
    {factors.map((factor, index) => {
      const groupA = data.filter(d => d[factor] === true);
      const groupB = data.filter(d => d[factor] === false);
      
      const nA = groupA.length;
      const nB = groupB.length;
      
      const pA = nA > 0 ? groupA.filter(d => d.interviewed || d.codeAssessmentRequested).length / nA : 0;
      const pB = nB > 0 ? groupB.filter(d => d.interviewed || d.codeAssessmentRequested).length / nB : 0;
      
      const delta = (pA - pB) * 100; 
      
      const varA = nA > 0 ? (Math.max(pA * (1 - pA), 0.05)) / nA : 0;
      const varB = nB > 0 ? (Math.max(pB * (1 - pB), 0.05)) / nB : 0;
      const moe = 1.96 * Math.sqrt(varA + varB) * 100;
      
      const isCalculable = nA > 0 && nB > 0;
      const displayDelta = isCalculable ? delta : 0;
      const displayMoe = isCalculable ? moe : 0;
      
      const isSignificant = isCalculable && (displayDelta - displayMoe > 0 || displayDelta + displayMoe < 0);
      
      let factorColor = 'gray';
      if (isSignificant) {
        factorColor = displayDelta > 0 ? 'lightgreen' : 'salmon';
      }
      const safeLeft = Math.max(-50, displayDelta - displayMoe);
      const safeRight = Math.min(50, displayDelta + displayMoe);
      const renderWidth = safeRight - safeLeft;
      
      return (
        <div key={factor} style={{ position: 'absolute', top: `${index * 30}px`, width: '100%', height: '30px' }}>
        <span style={{ position: 'absolute', left: '0', fontSize: '0.8em', top: '5px', opacity: isCalculable ? 1 : 0.4, zIndex: 3, background: 'var(--bg-secondary)', paddingRight: '4px' }}>
        {factorLabels[factor] || factor} {!isCalculable && "(N/A)"}
        </span>
        
        {isCalculable && (
          <>
          <div style={{
            position: 'absolute',
            left: `calc(50% + ${safeLeft}%)`,
            width: `${renderWidth}%`,
            height: '2px',
            background: factorColor,
            top: '14px',
            zIndex: 1
          }} />
          <div style={{
            position: 'absolute',
            left: `calc(50% + ${displayDelta}%)`,
            width: '12px',
            height: '12px',
            borderRadius: '50%',
            background: displayDelta > 0 ? 'lightgreen' : 'salmon',
            border: isSignificant ? '2px solid var(--bg-secondary)' : 'none',
            transform: 'translateX(-50%)',
            top: '9px',
            zIndex: 2
          }} title={`${factorLabels[factor] || factor} Impact: ${displayDelta.toFixed(1)}% ±${displayMoe.toFixed(1)}%`} />
          </>
        )}
        </div>
      );
    })}
    </div>
    </div>
    </div>
    </div>
  );
};