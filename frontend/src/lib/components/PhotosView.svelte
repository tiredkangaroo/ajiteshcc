<script lang="ts">
  import { isAdmin } from "$lib/stores/isAdmin";
  import type { Photo } from "../../types";

  let { photos } = $props();
  let imageInDialog: Photo | null = $state(null);
</script>

{#if imageInDialog}
  <div
    class="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50"
    onclick={() => {
      imageInDialog = null;
    }}
  >
    <img
      src={imageInDialog.photo_url}
      alt={imageInDialog.title}
      class="max-w-[90%] max-h-[90%] object-contain rounded-lg shadow-lg"
    />
  </div>
{/if}

<div class="flex flex-wrap w-full h-[90%] justify-center gap-4 overflow-y-auto">
  {#each photos as photo, index}
    <div class="relative group max-w-[400px] max-h-[300px]">
      <img
        src={photo.photo_url}
        alt={photo.title}
        class="w-full h-full object-cover rounded-lg"
        onclick={() => {
          imageInDialog = photo;
        }}
      />

      {#if photo.title || photo.comment || (photo.tags && photo.tags.length > 0)}
        <div
          class="mt-auto h-fit absolute inset-0 bg-black bg-opacity-60 opacity-0 group-hover:opacity-80 transition-opacity flex flex-col justify-center items-center text-center p-4"
        >
          {#if photo.title}
            <h3 class="text-white text-xl font-bold mb-2">
              {photo.title}
            </h3>
          {/if}
          {#if photo.comment}
            <p class="text-white text-sm mb-2">{photo.comment}</p>
          {/if}
          <div class="flex flex-wrap gap-2 justify-center">
            {#if photo.tags.length > 0}
              <span class="text-white font-semibold mr-2">tags:</span>
            {/if}
            {#each photo.tags as tag}
              <button
                class="px-2 py-1 bg-amber-300 text-sm font-light rounded"
                title={tag.comment}
                onclick={() => {
                  window.location.replace(`/tags/${tag.title}`);
                }}
              >
                {tag.title}
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/each}
</div>
