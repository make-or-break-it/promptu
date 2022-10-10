<script>
  let feed = [
    // SWITCHABLE - based on API results
    {
      message: "Hey, eating some green beans and ham today",
      user: "Cat in the Hat",
    },
    {
      message: "Running out of things to wear",
      user: "Mrs Piggy",
    },
    {
      message: "What year is it?",
      user: "alanparrish",
    },
    {
      message: "Had a nice, long run!",
      user: "Usain Bolt",
    },
    {
      message: "Ready to bake some lovely cakes!",
      user: "Mary Berry",
    },
  ];
</script>

<div class="feed">
  {#each feed as post}
    <div class="post">
      <p class="message">{post.message}</p>
      <p class="user">{post.user}</p>
    </div>
  {:else}
    <p class="overlay">No posts yet!</p>
  {/each}
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
  }
</style>
