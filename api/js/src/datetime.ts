const millisecond = 1;
const second = 1000 * millisecond;
const minute = 60 * second;
const hour = 60 * minute;
const day = 24 * hour;
const week = 7 * day;
const month = (365/12) * day;
const year = 365 * day;

export function relative(date: Date): string {
    const now = new Date(Date.now());
    const diff = now.valueOf() - date.valueOf();    // in milliseconds

    const years = Math.floor(diff / year);
    if (years > 1) {
        return `${years} years ago`;
    }
    if (years > 0) {
        return "1 year ago";
    }

    const months = Math.floor(diff / month);
    if (months > 1) {
        return `${months} months ago`;
    }
    if (months > 0) {
        return "1 month ago";
    }

    const weeks = Math.floor(diff / week);
    if (weeks > 1) {
        return `${weeks} weeks ago`;
    }
    if (weeks > 0) {
        return "1 week ago";
    }

    const days = Math.floor(diff / day);
    if (days > 1) {
        return `${days} days ago`;
    }
    if (days > 0) {
        return "1 day ago";
    }

    const hours = Math.floor(diff / hour);
    if (hours > 1) {
        return `${hours} hours ago`;
    }
    if (hours > 0) {
        return "1 hour ago";
    }

    const minutes = Math.floor(diff / minute);
    if (minutes > 1) {
        return `${minutes} minutes ago`;
    }
    if (minutes > 0) {
        return "1 minute ago";
    }

    const seconds = Math.floor(diff / second);
    if (seconds > 1) {
        return `${seconds} seconds ago`;
    }
    if (seconds > 0) {
        return "1 second ago";
    }
    return "Just now";
}

export function createDateTime(date: Date): HTMLSpanElement {
    const span = document.createElement("span");
    span.title = date.toLocaleString();
    span.textContent = relative(date);
    return span;
}
