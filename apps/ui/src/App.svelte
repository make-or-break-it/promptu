<script>
  import Feed from "./lib/Feed.svelte";
  import Header from "./lib/Header.svelte";
  import PostPrompt from "./lib/PostPrompt.svelte";
  import { posted } from "./stores/postManager"; // this needs to be replaced by an API
  import { notificationTime } from "./stores/postManager"; // this needs to be replaced by an API

  let timeToPost = new Date().getTime() > $notificationTime.getTime();
</script>

<main>
  <Header />
  {#if timeToPost && !$posted}
    <PostPrompt />
  {:else}
    <Feed />
  {/if}
</main>
