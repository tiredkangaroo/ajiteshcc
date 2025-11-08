<script>
  import { isAdmin } from "../lib/stores/isAdmin";
</script>

<div class="w-full h-full flex flex-col justify-center align-middle p-4">
  <div class="text-center">
    <h1 class="mt-4">{$isAdmin ? "hi, aji!" : "Ajitesh Kumar"}</h1>
    {#if $isAdmin}
      <a href="/admin">go to admin control panel</a>
    {/if}
    <img
      src="https://avatars.githubusercontent.com/u/81335306?v=4"
      alt="profile picture"
      class="w-48 h-48"
    />
  </div>
  <div class="flex flex-row justify-between items-center px-12">
    <h2><a href="/projects">projects 💻</a></h2>
    <h2><a href="/photography">photography 📷</a></h2>
    <h2><a href="/blog">blog 📝</a></h2>
  </div>
</div>
