<script>
  export let placeholder;
  export let value;
  export let validationMessage;
  export let validationRegex;
  
  let oldValue;
  let inputContent;

  let editContent = false;

  // using custom input handler instead of bind:textContent because two-way binding is not allowed with dynamic
  // dynamic coneteneditable attribute - see: https://stackoverflow.com/questions/57392773
  function handleInput(e) {
    inputContent = e.target.textContent
  }

  function handleDblClick(e) {
    oldValue = e.target.textContent
    editContent = true
  }

  function handleKeydown(e) {
    if (e.key == 'Enter') {
      e.preventDefault()

      const legalChars = validationRegex

      if (legalChars.test(inputContent)) {
        value = inputContent
      } else {
        alert(validationMessage)
        e.target.textContent = oldValue
      }

      editContent = false
    }
  }
</script>

<p class:editable="{editContent}" 
  contenteditable="{editContent}"
  on:dblclick|preventDefault={handleDblClick} 
  on:keydown={handleKeydown}
  on:input={handleInput}
  >
    {placeholder}
</p>

<style>
  p {
    text-align: center;
    padding: 1em 2em;
    font-size: medium;
    background-color: #004A8F;
    border-radius: 0.5em;
  }

  .editable {
    background-color: #00203D;
  }
</style>
