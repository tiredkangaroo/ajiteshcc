<script>
  import { isAdmin } from "../lib/stores/isAdmin";
  import { fetchBackend } from "../utils";
</script>

<div
  class="w-full h-full flex flex-col justify-center items-center p-8 bg-gradient-to-b from-amber-50 to-white"
>
  <div class="flex flex-row justify-between items-center w-full max-w-5xl p-6">
    <div class="flex flex-col gap-4 justify-center items-center">
      <h1 class="mt-2 text-4xl font-bold text-amber-900 drop-shadow-sm">
        {$isAdmin ? "hi, aji!" : "Ajitesh Kumar"}
      </h1>

      {#if $isAdmin}
        <div
          class="w-full flex flex-row justify-between text-blue-700 underline font-medium"
        >
          <a href="/admin">control panel</a>
          <button
            class="text-amber-700 underline"
            onclick={async () => {
              await fetchBackend("/api/v1/admin/logout", {
                method: "POST",
                credentials: "include",
              });
              window.location.reload();
            }}>log out</button
          >
        </div>
      {/if}

      <img
        src="https://avatars.githubusercontent.com/u/81335306?v=4"
        alt="tiredkangaroo profile"
        class="w-48 h-48 rounded-2xl"
      />
    </div>

    <div class="flex flex-col text-2xl gap-4 text-amber-800">
      <h2><a class="hover:underline" href="/photography">photography 📷</a></h2>
      <h2><a class="hover:underline" href="/blog">blog 📝</a></h2>
    </div>
  </div>
</div>
