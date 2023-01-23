<script>
  import { Icon } from "svelte-awesome";
  import users from "svelte-awesome/icons/users";
</script>

<button
  class:active={active === "activity"}
  on:click={() => (active = "activity")}
>
  <div class="indicator">
    <span
      class="indicator-item badge badge-xs badge-accent w-fit translate-x-8 md:translate-x-10"
      >new</span
    >
    <Icon data={users} scale={1.6} />
  </div>
  <p class="text-xs md:text-sm">Activity</p>
</button>
