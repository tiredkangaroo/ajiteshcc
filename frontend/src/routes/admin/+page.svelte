<script>
  import { page } from "$app/state";
  import { isAdmin } from "$lib/stores/isAdmin";
  import { fetchBackend } from "../../utils";

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
  let windowPage = $state("");
</script>

{#if $isAdmin}
  <div class="w-full h-full flex flex-row">
    <div class="min-w-[200px] w-1/5 bg-amber-200 py-1">
      <h1 class="text-center">hi aji!</h1>
      <div class="mt-6 text-center">
        <button
          class="text-xl px-2 py-4 underline"
          onclick={() => {
            windowPage = "blog";
          }}>Blog</button
        >
      </div>
    </div>
  </div>
{:else}
  <div class="w-full h-full flex justify-center items-center">
    <input
      type="text"
      placeholder="Enter code"
      class="border border-black px-4 py-2 rounded-2xl"
      bind:value={code}
    />
    <button
      class="ml-4 px-4 py-2 bg-amber-400 rounded-2xl"
      onclick={() => {
        submitCode();
      }}>submit</button
    >
  </div>
{/if}
