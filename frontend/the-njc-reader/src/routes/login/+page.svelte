<script>
  import Icon from "svelte-awesome";
  import userCircle from "svelte-awesome/icons/userCircle";
  import ellipsisH from "svelte-awesome/icons/ellipsisH";

  export let data;
  const returnURL = data.returnURL;

  export let form;
</script>

<div
  class="flex flex-col md:flex-row h gap-8 w-full justify-center items-center"
>
  <form method="POST">
    <div class="form-control">
      <label class="input-group">
        <span class="bg-base-100"
          ><Icon class="text-neutral" data={userCircle} scale={1.5} /></span
        >
        <input
          type="text"
          placeholder="Email"
          class="input input-bordered"
          name="email"
          required
        />
      </label>
    </div>
    <br />
    <div class="form-control">
      <label class="input-group">
        <span class="bg-base-100 "
          ><Icon class="text-neutral" data={ellipsisH} scale={1.8} /></span
        >
        <input
          type="password"
          placeholder="Password"
          class="input input-bordered"
          name="password"
          required
        />
      </label>
    </div>
    <br />
    <input type="hidden" name="returnURL" value={returnURL} hidden>
    <button class="btn btn-md bg-primary mx-3">Log in</button>
  </form>

  {#if form?.invalid}
    <p class="error">Username and password are required.</p>
  {/if}

  {#if form?.credentials}
    <p class="error">Invalid credentials.</p>
  {/if}

  <a href="/reset-password" class="link italic">Forgot password?</a>
</div>

<style>
  .h {
    min-height: 90vh;
  }
</style>
