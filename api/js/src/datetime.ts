export const millisecond = 1;
export const second = 1000 * millisecond;
export const minute = 60 * second;
export const hour = 60 * minute;
export const day = 24 * hour;
export const week = 7 * day;
export const month = (365 / 12) * day;
export const year = 365 * day;

export function relative(date: Date): string {
  const now = new Date(Date.now());
  if (date <= now) {
    return howLongAgo(now, date);
  }
  return howFarInTheFuture(now, date);
}

export function createDateTime(date: Date): HTMLSpanElement {
  const span = document.createElement("span");
  span.title = date.toLocaleString();
  span.textContent = relative(date);
  return span;
}

function howLongAgo(now: Date, date: Date): string {
  const diff = now.valueOf() - date.valueOf(); // in milliseconds

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

function howFarInTheFuture(now: Date, date: Date): string {
  const diff = date.valueOf() - now.valueOf(); // in milliseconds

  const years = Math.floor(diff / year);
  if (years > 1) {
    return `In ${years} years`;
  }
  if (years > 0) {
    return "Next year";
  }

  const months = Math.floor(diff / month);
  if (months > 1) {
    return `In ${months} months`;
  }
  if (months > 0) {
    return "Next month";
  }

  const weeks = Math.floor(diff / week);
  if (weeks > 1) {
    return `In ${weeks} weeks`;
  }
  if (weeks > 0) {
    return "Next week";
  }

  const days = Math.floor(diff / day);
  if (days > 1) {
    return `In ${days} days`;
  }
  if (days > 0) {
    return "Tomorrow";
  }

  const hours = Math.floor(diff / hour);
  if (hours > 1) {
    return `In ${hours} hours`;
  }
  if (hours > 0) {
    return "In 1 hour";
  }

  const minutes = Math.floor(diff / minute);
  if (minutes > 1) {
    return `In ${minutes} minutes`;
  }
  if (minutes > 0) {
    return "In 1 minute";
  }

  const seconds = Math.floor(diff / second);
  if (seconds > 1) {
    return `In ${seconds} seconds`;
  }
  if (seconds > 0) {
    return "In 1 second";
  }
  return "Now";
}

export function startOfDay(now: Date = new Date()): Date {
  const year = now.getFullYear();
  const month = now.getMonth();
  const date = now.getDate();
  return new Date(year, month, date);
}

export function endOfDay(now: Date = new Date()): Date {
  const year = now.getFullYear();
  const month = now.getMonth();
  const date = now.getDate();
  return new Date(year, month, date + 1);
}
