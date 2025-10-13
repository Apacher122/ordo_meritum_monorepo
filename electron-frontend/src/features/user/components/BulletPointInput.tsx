import React, { useEffect, useRef } from "react";

interface BulletPointInputProps {
  bullets: { text: string }[];
  onChange: (bullets: { text: string }[]) => void;
  label?: string;
}

export const BulletPointInput: React.FC<BulletPointInputProps> = ({
  bullets,
  onChange,
  label = "Bullet Points:",
}) => {
  const itemRefs = useRef<(HTMLTextAreaElement | null)[]>([]);

  useEffect(() => {
    itemRefs.current = itemRefs.current.slice(0, bullets.length);
  }, [bullets]);

  const handleTextChange = (index: number, value: string) => {
    const sanitizedValue = value.replace(/(\r\n|\n|\r)/gm, "");
    const newBullets = bullets.map((bullet, i) =>
      i === index ? { ...bullet, text: sanitizedValue } : bullet
    );
    onChange(newBullets);
  };

  const handleKeyDown = (
    e: React.KeyboardEvent<HTMLTextAreaElement>,
    index: number
  ) => {
    if (e.key === "Enter") {
      e.preventDefault();
      const newBullets = [
        ...bullets.slice(0, index + 1),
        { text: "" },
        ...bullets.slice(index + 1),
      ];
      onChange(newBullets);
      setTimeout(() => {
        itemRefs.current[index + 1]?.focus();
      }, 0);
    }

    if (e.key === "Backspace") {
      if (bullets[index].text.length === 0 && bullets.length > 1) {
        e.preventDefault();
        const newBullets = bullets.filter((_, i) => i !== index);
        onChange(newBullets);

        setTimeout(() => {
          const prevIndex = index > 0 ? index - 1 : 0;
          const prevRef = itemRefs.current[prevIndex];
          if (prevRef) {
            prevRef.focus();
            prevRef.setSelectionRange(prevRef.value.length, prevRef.value.length);
          }
        }, 0);
      }
    }
  };

  return (
    <div className="bullet-point-container">
      <label>{label}</label>
      {bullets.map((bullet, index) => (
        <div key={index} className="bullet-wrapper">
          <span>&bull;</span>
          <textarea
            ref={(el) => {
              itemRefs.current[index] = el;
            }}
            value={bullet.text}
            onChange={(e) => handleTextChange(index, e.target.value)}
            onKeyDown={(e) => handleKeyDown(e, index)}
            className="bullet-textarea"
            rows={1}
            placeholder="Add detail..."
          />
        </div>
      ))}
    </div>
  );
};