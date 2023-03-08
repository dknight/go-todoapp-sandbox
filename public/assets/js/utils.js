/**
 * @param {string}
 * @return {string}
 */
const chars = {
    "&": "&amp;",
    ">": "&gt;",
    "<": "&lt;",
    '"': "&quot;",
    "'": "&#39;",
    "`": "&#96;",
};

const re = new RegExp(Object.keys(chars).join("|"), "g");

const html = (str) => String(str).replace(re, (m) => chars[m]);

export {html};
