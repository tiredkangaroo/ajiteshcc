<script lang="ts">
  import { isAdmin } from "../lib/stores/isAdmin";
  import type { Track } from "../types";
  import { backendUrl, fetchBackend } from "../utils";

  let currentTrack: Track | null = $state(null);
  let getMusicError: string | null = $state(null);
  $effect(() => {
    async function fetchCurrentTrack() {
      const res = await fetchBackend("/api/v1/music");
      if (res.ok) {
        currentTrack = await res.json();
      } else {
        currentTrack = null;
        const errData = await res.json();
        getMusicError = errData.error || "unknown error";
      }
    }
    fetchCurrentTrack();
  });
</script>

<div
  class="w-full h-full flex flex-col justify-center items-center p-8 bg-gradient-to-b from-amber-50 to-white"
>
  <div class="flex flex-row justify-between items-center w-full max-w-5xl p-6">
    <div class="flex flex-col gap-4 justify-center items-center">
      <h1 class="mt-2 text-4xl font-bold text-amber-900">
        {$isAdmin ? "hi, aji!" : "hi, i'm ajitesh!"}
      </h1>

      {#if $isAdmin}
        <div
          class="w-full flex flex-row justify-between text-blue-700 underline font-medium"
        >
          <a href="/admin">control panel</a>
          <button
            class="underline"
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
      <h2>
        <a
          class="text-amber-700 underline visited:text-amber-800"
          href="/photography">photography 📷</a
        >
      </h2>
      <h2>
        <a class="text-amber-700 underline visited:text-amber-800" href="/blog"
          >blog 📝</a
        >
      </h2>
    </div>
  </div>
  {#if currentTrack}
    <div
      class="bg-amber-900 text-amber-100 px-8 py-3 rounded-xl flex flex-col gap-3 mt-6"
    >
      <div class="w-full items-center flex flex-row gap-2 justify-center">
        <p class="text-center">what's aji playing?</p>
        {#if $isAdmin}
          <p class="text-center text-sm">
            <a href={backendUrl + "/api/v1/music/logout"} class="underline"
              >logout</a
            >
          </p>
        {/if}
      </div>
      <div class="flex flex-row items-center gap-9">
        <div class="relative w-16 h-16">
          <img
            src={currentTrack.cover_url}
            alt="album cover"
            class="w-full h-full rounded-full animate-[spin_5s_linear_infinite] opacity-90"
          />
          <div
            class="absolute inset-0 rounded-full bg-gradient-conic from-transparent via-white/40 to-transparent animate-[spin_5s_linear_infinite] pointer-events-none"
          ></div>
          <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 rounded-full bg-gray-900 border-2 border-gray-700"
          ></div>
        </div>
        <div class="flex flex-col">
          <p class="text-lg font-semibold text-amber-100">
            {currentTrack.name}
          </p>
          <p class="text-xs text-amber-200">{currentTrack.artists}</p>
        </div>
      </div>
    </div>
  {:else if getMusicError && $isAdmin}
    <div
      class="bg-red-900 text-amber-100 px-8 py-3 rounded-xl flex flex-col gap-3 mt-6"
    >
      <p class="text-center">error: {getMusicError}</p>
      <p>
        perhaps... <a
          href={backendUrl + "/api/v1/music/login"}
          class="text-gray-100 underline">log in to spotify</a
        >
        or
        <a
          href={backendUrl + "/api/v1/music/logout"}
          class="text-gray-100 underline">log out</a
        >?
      </p>
    </div>
  {/if}
</div>
