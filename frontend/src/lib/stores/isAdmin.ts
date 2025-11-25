import { writable, type Writable } from "svelte/store";

export const isAdmin: Writable<boolean | undefined> = writable(undefined);
