<script lang="ts">
  import PhotosView from "$lib/components/PhotosView.svelte";
  import type { Photo } from "../../types";
  import { fetchBackend } from "../../utils";

  let photos: Photo[] | null = $state(null);

  $effect(() => {
    async function loadPhotos() {
      const res = await fetchBackend("/api/v1/photos");
      if (res.ok) {
        const data = await res.json();
        if (data === null) {
          photos = [];
          return;
        }
        photos = data;
      } else {
        console.error("failed to load photos", await res.text());
      }
    }
    loadPhotos();
  });
</script>

{#if photos === null}
  <div class="w-full h-full justify-center items-center flex">
    Loading photos...
  </div>
{:else if photos.length === 0}
  <div class="w-full h-full flex flex-col">
    {@render title()}
    <div class="h-full justify-center items-center flex">No photos.</div>
  </div>
{:else}
  <div class="h-full overflow-hidden">
    {@render title()}
    <div class="flex flex-row w-full h-full overflow-hidden">
      <PhotosView {photos} />
    </div>
  </div>
{/if}

{#snippet title()}
  <div class="w-full relative flex items-center justify-center mt-3 mb-4">
    <button
      onclick={() => {
        window.location.replace("/");
      }}
      class="absolute left-4 px-4 py-1 bg-amber-400">back</button
    >
    <h1 class="text-xl font-semibold">photography 📷</h1>
  </div>
{/snippet}
