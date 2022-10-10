<script>
  import Feed from "./lib/Feed.svelte";
  import Header from "./lib/Header.svelte";
  import PostPrompt from "./lib/PostPrompt.svelte";

  let timeToPost = false; // SWITCHABLE - based on API results
  let posted = false; // SWITCHABLE - based on API results
</script>

<main>
  <Header />
  {#if timeToPost && !posted}
    <PostPrompt />
  {:else}
    <Feed />
  {/if}
</main>
