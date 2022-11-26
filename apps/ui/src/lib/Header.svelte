<script>
  import { username } from "../stores/username";
  import EditableButton from "./EditableButton.svelte";

  const PromptuLogo = "Promptu";
</script>

<ul>
  <li class="left-aligned">
    <a href="/">{PromptuLogo}</a>
  </li>
  <li class="right-aligned">
    <EditableButton 
      bind:value={$username} 
      validationMessage="Username must be at least one character long and contain alphanumeric characters, - or _"
      validationRegex={/^[\w\d\-_]+$/}
      placeholder={$username ? $username : "Set username"} />
  </li>
</ul>

<style>
  ul {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    list-style-type: none;
    margin: 0;
    padding: 0.5em 2em;
    border-bottom: solid 3px whitesmoke;
    font-size: 1.5em;
    font-weight: bold;
    background-color: #003566;
    box-shadow: 0 0 1em 0 #1c2541;
  }

  a {
    line-height: 100px;
    text-decoration: none;
    display: block;
    text-align: center;
    color: whitesmoke;
  }

  li {
    display: flex;
    flex-direction: row;
    flex-grow: 1;
    align-items: center;
  }

  .left-aligned {
    justify-content: flex-start;
  }

  .right-aligned {
    justify-content: flex-end;
  }
</style>
