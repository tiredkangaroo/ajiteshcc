<script lang="ts">
  import "../app.css";
  import { isAdmin } from "../lib/stores/isAdmin";
  import { fetchBackend } from "../utils";

  let { children } = $props();

  $effect(() => {
    async function checkAdmin() {
      const res = await fetchBackend("/api/v1/admin", {
        method: "GET",
        credentials: "include",
      });
      const data = await res.json();
      console.log(data);
      isAdmin.set(data.is_admin);
    }
    checkAdmin();
  });
</script>

{#if $isAdmin === undefined}
  <div></div>
{:else}
  {@render children()}
{/if}
