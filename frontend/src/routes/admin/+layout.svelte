<script lang="ts">
  import { page } from "$app/state";
  import { isAdmin } from "$lib/stores/isAdmin";
  import { fetchBackend } from "../../utils";

  let { children } = $props();
  let code = $state("");

  function submitCode() {
    console.log(code);
    fetchBackend("/api/v1/admin", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({
        totp: code,
      }),
    }).then(async (res) => {
      if (res.ok) {
        isAdmin.set(true);
      } else {
        alert("invalid code");
      }
    });
  }
</script>

{#if $isAdmin}
  <div class="w-full h-full flex flex-row">
    <div
      class="min-w-[220px] w-min bg-amber-100 py-4 px-3 shadow-xl border-r border-amber-300 rounded-r-2xl overflow-none"
    >
      <h1 class="text-center text-2xl font-bold text-amber-900">hi aji!</h1>

      <div class="mt-8 flex flex-col space-y-4">
        <button
          class="text-xl px-3 py-2 rounded-xl bg-white shadow hover:bg-amber-50 transition"
          onclick={() => {
            window.location.assign("/admin/blog");
          }}
          style={page.url.pathname.startsWith("/admin/blog")
            ? "background-color: oklch(47.3% 0.137 46.201); color: #fef3c6;"
            : "background-color: white; color: oklch(47.3% 0.137 46.201);"}
        >
          blog
        </button>

        <button
          class="text-xl px-3 py-2 rounded-xl shadow hover:bg-amber-50 transition"
          style={page.url.pathname.startsWith("/admin/photos")
            ? "background-color: oklch(47.3% 0.137 46.201); color: #fef3c6;"
            : "background-color: white; color: oklch(47.3% 0.137 46.201);"}
          onclick={() => {
            window.location.assign("/admin/photos");
          }}
        >
          photos
        </button>
      </div>
    </div>
    <div class="w-full h-full overflow-auto">
      {@render children()}
    </div>
  </div>
{:else}
  <div class="w-full h-full flex justify-center items-center">
    <input
      type="text"
      placeholder="Enter code"
      class="border border-black px-4 py-2"
      bind:value={code}
    />
    <button
      class="ml-4 px-4 py-2 bg-amber-400"
      onclick={() => {
        submitCode();
      }}>submit</button
    >
  </div>
{/if}
