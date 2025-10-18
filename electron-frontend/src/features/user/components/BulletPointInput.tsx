import React, { useEffect, useMemo, useRef } from "react";

interface BulletPoint {
  text: string;
  id?: string;
}

interface BulletPointInputProps {
  bullets: BulletPoint[];
  onChange: (bullets: BulletPoint[]) => void;
  label?: string;
}

const generateId = () => Math.random().toString(36).slice(2, 11);

export const BulletPointInput: React.FC<BulletPointInputProps> = ({
  bullets,
  onChange,
  label = "Bullet Points:",
}) => {
  const itemRefs = useRef<(HTMLTextAreaElement | null)[]>([]);
  
  const keyedBullets = useMemo(() => {
    return bullets.map(bullet => ({
      text: bullet.text || "",
      id: bullet.id || generateId()
    }));
  }, [bullets]);
  
  
  const adjustHeight = (el: HTMLTextAreaElement | null) => {
    if (el) {
      el.style.height = 'auto';
      el.style.height = `${el.scrollHeight}px`;
    }
  };
  
  useEffect(() => {
    for (const item of itemRefs.current) {
      adjustHeight(item);
    }
  }, [keyedBullets]);
  
  useEffect(() => {
    itemRefs.current = itemRefs.current.slice(0, bullets.length);
  }, [bullets]);
  
  const handleTextChange = (index: number, value: string) => {
    const stableId = keyedBullets[index]?.id;
    
    const newBullets = bullets.map((bullet, i) => {
        const existingId = keyedBullets[i]?.id;
        if (i === index) {
            return { text: value, id: stableId }; 
        }
        return { text: bullet.text, id: existingId }; 
    });
    onChange(newBullets);
    adjustHeight(itemRefs.current[index]);

};

const handleKeyDown = (
  e: React.KeyboardEvent<HTMLTextAreaElement>,
  index: number
) => {
  if (e.key === "Enter") {
    e.preventDefault();
    const newBullet: BulletPoint = { text: "", id: generateId() };
    const newBullets = [
      ...bullets.slice(0, index + 1),
      newBullet,
      ...bullets.slice(index + 1),
    ];
    onChange(newBullets);
    setTimeout(() => {
      adjustHeight(itemRefs.current[index + 1]);
      itemRefs.current[index + 1]?.focus();
    }, 0);
  }
  
  if (e.key === "Backspace") {
    const currentTextLength = String(bullets[index].text || '').length;
    
    if (currentTextLength === 0 && bullets.length > 1) {
      e.preventDefault();
      const newBullets = bullets.filter((_, i) => i !== index);
      onChange(newBullets);
      
      setTimeout(() => {
        const prevIndex = index > 0 ? index - 1 : 0;
        const prevRef = itemRefs.current[prevIndex];
        if (prevRef) {
          prevRef.focus();
          prevRef.setSelectionRange(prevRef.value.length, prevRef.value.length);
          adjustHeight(prevRef);
        }
      }, 0);
    }
  }
};

const handlePaste = (e: React.ClipboardEvent<HTMLTextAreaElement>, index: number) => {
    const pasteText = e.clipboardData.getData('text');
    
    if (pasteText.includes('\n') || pasteText.includes('\r')) {
      e.preventDefault();
      
      const lines = pasteText
        .split(/[\r\n]+/)
        .map(line => line.trim())
        .filter(line => line.length > 0);
      
      if (lines.length > 0) {
        const firstLine = lines[0];
        
        const newBulletPoints = lines.slice(1).map(text => ({ text, id: generateId() }));
        
        const newBullets = bullets.map((bullet, i) => {
          const existingId = keyedBullets[i]?.id; 
          if (i === index) {
            const stableId = keyedBullets[index]?.id;
            return { text: firstLine, id: stableId }; 
          }
          return { text: bullet.text, id: existingId }; 
        });

        const updatedBullets = [
          ...newBullets.slice(0, index + 1),
          ...newBulletPoints,
          ...newBullets.slice(index + 1),
        ];

        onChange(updatedBullets);
        
        setTimeout(() => {
            const nextIndex = index + 1;
            adjustHeight(itemRefs.current[index]);
            itemRefs.current[nextIndex]?.focus();
        }, 0);
      }
    }
  };

  return (
    <div className="bullet-point-container">
      <label>{label}</label>
      {keyedBullets.map((bullet, index) => (
        <div key={bullet.id} className="bullet-wrapper">
          <span>&bull;</span>
          <textarea
            ref={(el) => {
              itemRefs.current[index] = el;
            }}
            value={bullet.text || ""} 
            onChange={(e) => handleTextChange(index, e.target.value)}
            onKeyDown={(e) => handleKeyDown(e, index)}
            onPaste={(e) => handlePaste(e, index)}
            className="bullet-textarea"
            rows={1}
            placeholder="Add detail..."
          />
        </div>
      ))}
    </div>
  );
}