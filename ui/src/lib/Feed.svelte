<script>
  const feedEndpoint = 'http://localhost:8080/feed'

  const getFeed = (async () => {
    let response = await fetch(feedEndpoint)
    let feedContent = await response.json()

    if ("posts" in feedContent) {
      return feedContent["posts"]
    } else {
      return []
    }
  })()
</script>

<div class="feed">
  {#await getFeed}
    <p class="overlay">No posts yet!</p>
  {:then feed}
    {#each feed as post}
      <div class="post">
        <p class="message">{post.message}</p>
        <p class="user">{post.user}</p>
      </div>
    {:else}
      <p class="overlay">No posts yet!</p>
    {/each}
  {:catch err}
    <p class="overlay error">Error loading posts - computer says: <span class="error-code">{err.message}</span></p>
    <p class="overlay error">Try reloading the page</p>
  {/await}
</div>

<style>
  p {
    font-size: 1.5em;
    font-weight: bold;
  }

  .feed {
    display: flex;
    justify-items: flex-start;
    flex-direction: column;
    padding: 2em 7em;
  }

  .post {
    color: black;
    background-color: rgba(255, 255, 255, 0.4);
    padding: 1em 1.5em;
    margin: 1em;
    border-radius: 0.5em;
  }

  .user {
    font-size: 0.7em;
    color: grey;
  }

  .overlay {
    text-align: center;
    color: #2E282A;
  }

  .error {
    color: #84110B;
  } 

  .error-code {
    font-weight: normal;
    background: #574C50;
    color: white;
    padding: 0.3em;
    border-radius: 0.2em;
  }
</style>
