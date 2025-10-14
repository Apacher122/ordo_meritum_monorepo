import { Experience, Project } from "../../types";

export const isBulletPointValid = (bp: { text: string }): boolean => {
    return bp.text.trim().length > 0; 
};

export const hasValidBulletPoints = (item: Experience | Project): boolean => {
    return (item.bulletPoints?.length ?? 0) > 0;
};

export const cleanBulletPoints = <T extends Experience | Project>(items: T[] | undefined): T[] => {
    if (!items) return [];
    
    const itemsWithCleanedBullets = items.map(item => ({
        ...item,
        bulletPoints: item.bulletPoints?.filter(isBulletPointValid)
    }));
    return itemsWithCleanedBullets.filter(hasValidBulletPoints) as T[];
};