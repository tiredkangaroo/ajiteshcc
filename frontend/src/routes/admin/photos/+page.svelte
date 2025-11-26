<script lang="ts">
  import PhotosView from "$lib/components/PhotosView.svelte";
  import type { Photo, Object } from "../../../types";
  import { fetchBackend } from "../../../utils";

  let photos = $state<Array<Photo> | undefined>(undefined);
  let objects = $state<Array<Object> | undefined>(undefined);
  let privatePhotos: Array<Object> | undefined = $derived.by(() => {
    if (objects == undefined || photos == undefined) {
      return undefined;
    }
    // check if object photo_url is in photos
    // if not, then it is a private photo
    const photoUrls = new Set(photos.map((photo) => photo.photo_url));
    return objects.filter((obj) => !photoUrls.has(obj.public_url)); // filter out public photos
  });
  $inspect(
    "photos",
    photos,
    "objects",
    objects,
    "privatePhotos",
    privatePhotos
  );

  $effect(() => {
    async function loadObjectsAndPhotos() {
      const response = await fetchBackend("/api/v1/objects");
      objects = await response.json();
      const photosResponse = await fetchBackend("/api/v1/photos");
      photos = await photosResponse.json();
    }
    loadObjectsAndPhotos();
  });
</script>

<div class="px-2 py-2 overflow-none">
  <h1 class="">photos</h1>
  <div class="overflow-none px-4 mt-4">
    <h3 class="text-xl font-semibold">public photos</h3>
    <div class="overflow-y-auto max-h-[40vh]">
      <PhotosView {photos} />
    </div>
    <h3 class="text-xl font-semibold mt-8">private photos</h3>
    <div class="overflow-y-auto max-h-[40vh]">
      <div
        class="flex flex-wrap w-full h-[90%] justify-center gap-4 overflow-y-auto"
      >
        {#each photos as photo, index}
          <div class="relative group max-w-[400px] max-h-[300px]">
            <img
              src={photo.photo_url}
              alt={photo.title}
              class="w-full h-full object-cover rounded-lg"
            />
          </div>
        {/each}
      </div>
    </div>
  </div>
</div>
