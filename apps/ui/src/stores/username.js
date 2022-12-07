import { writable } from "svelte/store";
import { posted } from "./postManager";

// Ref: https://dev.to/danawoodman/svelte-quick-tip-connect-a-store-to-local-storage-4idi
const localStoreUsername = localStorage.username

export const username = writable(localStoreUsername || '')

// Don't need to call unsubscribe as long as it's referenced through the $ prefix: https://svelte.dev/tutorial/auto-subscriptions
username.subscribe((u) => { 
    localStorage.username = u
    localStorage.getItem(`${u}Posted`) || posted.set(false)
})