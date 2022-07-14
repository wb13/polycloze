const canvas = document.createElement("canvas");

export function getFont(ele: Element): string {
    const style = getComputedStyle(ele);
    const weight = style.getPropertyValue("font-weight") || "normal";
    const size = style.getPropertyValue("font-size") || "12px";
    const family = style.getPropertyValue("font-family") || "sans-serif";
    return `${weight} ${size} ${family}`;
}

export function getWidth(font: string, text: string): string {
    const context = canvas.getContext("2d") as CanvasRenderingContext2D;
    context.font = font;
    const metrics = context.measureText(text);
    return `${metrics.width}px`;
}
