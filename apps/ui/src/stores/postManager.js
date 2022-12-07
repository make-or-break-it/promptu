import { readable, writable } from "svelte/store";

const localStoreUsername = localStorage.getItem('username')

const localStorePosted = localStorage.getItem(`${localStoreUsername}Posted`)
export const posted = writable(JSON.parse(localStorePosted) || false)
posted.subscribe((v) => { 
    const localStoreUsername = localStorage.getItem('username')
    v && localStorage.setItem(`${localStoreUsername}Posted`, JSON.stringify(v)) 
})

let currTime = new Date();
export const notificationTime = readable(new Date(currTime.getFullYear(), currTime.getMonth(), currTime.getDate(), 0, 0)); // for more authentic experience, turn into API)