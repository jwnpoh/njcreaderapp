<script>
  export let data;
  export let form;

  const email = data.user.email;
  const pmclass = data.user.class;
  let userName = data.user.display_name;

  let save = false;
  let newPassword;
  const saveChanges = () => {
    if (newPassword === "" && userName === data.user.display_name) {
      save = false;
    } else {
      save = true;
    }
  };

  let password = true;
  const passwordEnterd = () => {
    password = false;
  };
</script>

<form method="POST" action="?/updateUser">
  <div class="flex pt-10 justify-center items-center">
    <div class="card md:w-7/12 bg-base-100 shadow-sm shadow-neutral">
      <div class="card-body items-center text-center">
        <h1 class="card-title pb-2 text-2xl">User Profile</h1>
        <h2 class="text-lg font-semibold">Email</h2>
        <p class="text-neutral opacity-80">{email}</p>
        <h2 class="text-lg font-semibold">Class</h2>
        <p class="text-neutral opacity-80">{pmclass}</p>
        <h2 class="text-lg font-semibold">Display name</h2>
        <input
          name="display_name"
          type="text"
          class="input w-full max-w-md text-center bg-neutral bg-opacity-5"
          bind:value={userName}
          on:input={saveChanges}
        />
        <h2 class="pt-10 text-md font-semibold">Change password?</h2>
        <input
          name="new_password"
          type="password"
          placeholder="Enter new password"
          class="input w-full max-w-md text-center bg-neutral bg-opacity-5"
          bind:value={newPassword}
          on:input={saveChanges}
        />
        {#if save}
          <h2 class="pt-3 text-md font-semibold">Confirm new password</h2>
          <input
            required
            name="confirm"
            type="password"
            placeholder="Confirm new password"
            class="input w-full max-w-md text-center bg-neutral bg-opacity-5"
          />
          <h2 class="pt-3 text-md font-extrabold">Current password</h2>
          <input
            required
            name="old_password"
            type="password"
            placeholder="Enter current password to confirm changes."
            class="input w-full max-w-md text-center bg-primary bg-opacity-10"
            class:bg-primary={password}
            on:input={passwordEnterd}
          />
        {/if}
        <div class="card-actions justify-end pt-2">
          <button class="btn btn-primary">Save changes</button>
        </div>
        {#if form?.sent}
          <div class="alert alert-success max-w-fit place-self-center">
            <span>Changes saved successfully.</span>
          </div>
        {/if}
        {#if form?.failed}
          <div class="alert alert-error max-w-fit place-self-center">
            <span class="text-center">{form?.message}</span>
          </div>
        {/if}
      </div>
    </div>
  </div>
</form>
