// Creates an event listener that listens to clicks outside of an element.
// Usage: document.addEventListener("click", onClickOutside(elem, callback));
// The callback function takes the original target of the click event as input.
export function onClickOutside(elem: Element, callback: (target: EventTarget) => void): (event: Event) => void {
    return (event: Event) => {
        const target = event.target;
        if (!(target instanceof Node)) {
            return;
        }
        if (elem != target && target.contains(elem)) {
            callback(target);
        }
    };
}
