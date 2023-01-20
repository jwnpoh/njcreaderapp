<script>
  export let form;
</script>

<div
  class="grid pt-10 gap-2 w-full justify-center items-center"
>
  <p class="py-3">Enter the email address that was registered with your account.</p>
  <form method="POST">
    <div class="form-control">
      <label class="input-group">
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
    <button class="btn btn-md bg-primary ">Reset password</button>
  </form>

  {#if form?.invalid}
    <p class="error">{form.message}</p>
  {/if}

  {#if form?.success}
    <div class="alert alert-success max-w-fit place-self-center">
      <span>Password reset completed. Please check your inbox for the new temporary password.</span>
    </div>
  {/if}
</div>
