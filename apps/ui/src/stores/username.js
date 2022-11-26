import { writable } from "svelte/store";


// TODO: follow instructions in https://dev.to/danawoodman/svelte-quick-tip-connect-a-store-to-local-storage-4idi
const localStoreUsername = localStorage.username

export const username = writable(localStoreUsername || '')

// TODO: do we need to unsubscribe from this? Maybe need to move this to App.svelte
username.subscribe((v) => { localStorage.username = v})