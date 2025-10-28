import { Experience, Project } from "../../types";

/**
 * Checks if a bullet point object is valid by ensuring its `text` property
 * is a non-empty, trimmed string.
 * @param {object} bp - The bullet point object.
 * @returns {boolean} True if the bullet point is valid, false otherwise.
 */
export const isBulletPointValid = (bp: { text: string }): boolean => {
    return bp.text.trim().length > 0; 
};

/**
 * Checks if an Experience or Project item has at least one valid bullet point.
 * @param {Experience | Project} item - The experience or project item.
 * @returns {boolean} True if the item has valid bullet points.
 */
export const hasValidBulletPoints = (item: Experience | Project): boolean => {
    return (item.bulletPoints?.length ?? 0) > 0;
};

/**
 * Filters an array of Experience or Project items, removing any items that do not
 * have valid bullet points, and also removes invalid bullet points from within the items.
 * @template T
 * @param {T[] | undefined} items - The array of items to clean.
 * @returns {T[]} A new array containing only items with valid bullet points.
 */
export const cleanBulletPoints = <T extends Experience | Project>(items: T[] | undefined): T[] => {
    if (!items) return [];
    
    const itemsWithCleanedBullets = items.map(item => ({
        ...item,
        bulletPoints: item.bulletPoints?.filter(isBulletPointValid)
    }));
    return itemsWithCleanedBullets.filter(hasValidBulletPoints) as T[];
};