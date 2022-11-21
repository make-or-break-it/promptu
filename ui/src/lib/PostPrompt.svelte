<script>
  import { username } from "../stores/username";
  import { posted } from "../stores/postManager";

  const env = import.meta.env

  const postEndpoint = `${env.PUBLIC_PROMPTU_API_HOST}/post`
  const prompt = "Living or dead, who would you spend 3 hours with and why?";

  let answer = "";

  async function handleSubmit(e) {
    if (username && answer) {
      let response = await fetch(postEndpoint, {
        method: 'POST',
        body: JSON.stringify({
          'user': $username,
          'message': answer
        }),
         headers: {
          'Content-Type': 'application/json'
        }
      })

      let ok = await response.ok

      if (ok) {
        $posted = true
      }
    }
  }
</script>

<div>
  <p class="message">It's time to post! Today's prompt is:</p>
  <p class="prompt">{prompt}</p>

  <form on:submit|preventDefault={handleSubmit}>
    <input type="text" bind:value={answer} placeholder="Answer" />
    <input type="submit" value="Post!" />
  </form>
</div>

<style>
  div {
    width: 100%;
    height: 72vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 10em 0 0 0;
    background-color: rgba(0, 0, 0, 0.3);
  }

  .message {
    font-size: 1.5em;
  }

  .prompt {
    font-size: 2em;
    font-weight: bold;
  }

  form {
    margin: 2em 0;
  }

  input[type="text"] {
    padding: 0.5em;
    width: 30em;
    height: 2em;
    font-size: 1em;
    word-wrap: break-word;
    word-break: break-all;
  }

  input[type="submit"] {
    padding: 0.5em;
    height: 3.2em;
    font-size: 1em;
    font-weight: bold;
    background-color: #06d6a0;
    border-radius: 0.2em;
    border: 0;
  }
</style>
