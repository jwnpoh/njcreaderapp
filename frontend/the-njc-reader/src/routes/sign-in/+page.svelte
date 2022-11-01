<script>
  import Icon from "svelte-awesome";
  import userCircle from "svelte-awesome/icons/userCircle";
  import ellipsisH from "svelte-awesome/icons/ellipsisH";
</script>

<div
  class="flex flex-col md:flex-row h gap-8 w-full justify-center items-center"
>
  <div class="form-control">
    <label class="input-group">
      <span><Icon data={userCircle} scale={1.2} /></span>
      <input
        type="text"
        placeholder="info@site.com"
        class="input input-bordered"
        name="email"
      />
    </label>
  </div>
  <div class="form-control">
    <label class="input-group">
      <span><Icon data={ellipsisH} scale={1.5} /></span>
      <input
        type="password"
        placeholder="Password"
        class="input input-bordered"
        name="password"
      />
    </label>
  </div>
</div>

<style>
  .h {
    min-height: 90vh;
  }
</style>
