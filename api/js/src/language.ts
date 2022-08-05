// Language info in localStorage.

export type Language = {
  code: string
  name: string
  bcp47: string
}

export function getL1(): Language {
    return {
        code: localStorage.getItem("l1.code") || "eng",
        name: localStorage.getItem("l1.name") || "English",
        bcp47: localStorage.getItem("l1.bcp47") || "en",
    };
}

export function getL2(): Language {
    return {
        code: localStorage.getItem("l2.code") || "spa",
        name: localStorage.getItem("l2.name") || "Spanish",
        bcp47: localStorage.getItem("l2.bcp47") || "es",
    };
}

function setLanguage(prefix: "l1" | "l2", language: Language) {
    localStorage.setItem(`${prefix}.code`, language.code);
    localStorage.setItem(`${prefix}.name`, language.name);
    localStorage.setItem(`${prefix}.bcp47`, language.bcp47);
}

function swapL1L2() {
    const l1 = getL1();
    const l2 = getL2();
    setLanguage("l1", l2);
    setLanguage("l2", l1);
}

// Swaps L1 and L2 if needed to make sure that L1 != L2.
export function setL1(language: Language) {
    if (language.code !== getL2().code) {
        setLanguage("l1", language);
    } else {
        swapL1L2();
    }
}

// Swaps L1 and L2 if needed to make sure that L1 != L2.
export function setL2(language: Language) {
    if (language.code !== getL1().code) {
        setLanguage("l2", language);
    } else {
        swapL1L2();
    }
}
