import { writable } from "svelte/store";

const isBrowser = typeof window !== "undefined";

const stored = isBrowser ? sessionStorage.getItem("articleQueue") : null;
const initial = stored
  ? JSON.parse(stored).map(e => ({ ...e, date: new Date(e.date) }))
  : [];

const articleQueue = writable(initial);

articleQueue.subscribe(value => {
  if (isBrowser) {
    sessionStorage.setItem("articleQueue", JSON.stringify(
      value.map(e => ({ ...e, date: e.date.toISOString() }))
    ));
  }
});

export default articleQueue;
