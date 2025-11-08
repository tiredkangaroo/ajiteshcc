<script>
  import { isAdmin } from "../lib/stores/isAdmin";
  import { fetchBackend } from "../utils";
</script>

<div class="w-full h-full flex flex-col justify-center align-middle p-4">
  <div class="flex flex-row justify-between items-center px-12">
    <div class="flex flex-col gap-2 justify-center items-center">
      <h1 class="mt-4">{$isAdmin ? "hi, aji!" : "Ajitesh Kumar"}</h1>
      {#if $isAdmin}
        <div class="flex flex-row text-blue text-underline justify-between">
          <a href="/admin">control panel</a>
          <button
            class="text-amber-700 underline"
            onclick={async () => {
              await fetchBackend("/api/v1/admin/logout", {
                method: "POST",
                credentials: "include",
              });
              window.location.reload();
            }}
            >log out
          </button>
        </div>
      {/if}
      <img
        src="https://avatars.githubusercontent.com/u/81335306?v=4"
        alt="tiredkangaroo profile"
        class="w-48 h-48"
      />
    </div>
    <!-- <h2><a href="/projects">projects 💻</a></h2> -->
    <h2><a href="/photography">photography 📷</a></h2>
    <h2><a href="/blog">blog 📝</a></h2>
  </div>
</div>
