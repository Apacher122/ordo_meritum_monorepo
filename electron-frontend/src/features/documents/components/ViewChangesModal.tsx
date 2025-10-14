import '@/assets/styles/Components/UI/ViewChangesModal.css';

import React, { useEffect } from 'react';

import { ResumeChanges } from '../types/index';

interface ViewChangesModalProps {
    isOpen: boolean;
    onClose: () => void;
    changes: ResumeChanges;
}

export const ViewChangesModal: React.FC<ViewChangesModalProps> = ({ isOpen, onClose, changes }) => {
    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === 'Escape') {
                onClose();
            }
        };

        if (isOpen) {
            document.addEventListener('keydown', handleKeyDown);
        }

        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [isOpen, onClose]);

    if (!isOpen || !changes) return null;
    
    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={e => e.stopPropagation()}>
                <button className="modal-close-button" onClick={onClose}>&times;</button>
                <h2>Resume Changes & Justifications</h2>

                <div className="changes-section">
                    <h3>Summary</h3>
                    {changes.summary.map((body, i) => (
                        <div key={`summary-change-${i}-${body.sentence.slice(0, 10)}`} className="change-card">
                            <ul>
                                <li key={`summary-item-${body.sentence.slice(0, 15)}-${i}`} className={body.is_new_suggestion ? 'highlight' : ''}>
                                    <p>{body.sentence}</p>
                                    <small>Justification: {body.justification_for_change}</small>
                                </li>
                            </ul>
                        </div>
                    ))}
                </div>
                
                <div className="changes-section">
                    <h3>Experience</h3>
                    {changes.experiences.map((exp, i) => (
                        <div key={`exp-${i}-${exp.company}`} className="change-card">
                            <h4>{exp.position} at {exp.company}</h4>
                            <ul>
                                {exp.bulletPoints.map((desc, j) => (
                                    <li key={`exp-bullet-${j}-${desc.text.slice(0, 15)}`} className={desc.is_new_suggestion ? 'highlight' : ''}>
                                        <p>{desc.text}</p>
                                        <small>Justification: {desc.justification_for_change}</small>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>

                <div className="changes-section">
                    <h3>Projects</h3>
                    {changes.projects.map((proj, i) => (
                        <div key={`proj-${i}-${proj.name}`} className="change-card">
                            <h4>{proj.role} at {proj.name}</h4>
                            <ul>
                                {proj.bulletPoints.map((desc, j) => (
                                    <li key={`proj-bullet-${j}-${desc.text.slice(0, 15)}`} className={desc.is_new_suggestion ? 'highlight' : ''}>
                                        <p>{desc.text}</p>
                                        <small>Justification: {desc.justification_for_change}</small>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>

                 <div className="changes-section">
                    <h3>Skills</h3>
                    {changes.skills.map((skill, i) => (
                         <div key={`skill-category-${skill.category}-${i}`} className="change-card">
                            <h4>{skill.category}</h4>
                            <p><small>Justification: {skill.justification_for_changes}</small></p>
                            <ul>{skill.skill.map((s, j) => <li key={`skill-item-${s}-${j}`}>{s}</li>)}</ul>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

