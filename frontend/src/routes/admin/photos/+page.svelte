<script lang="ts">
  import PhotosView from "$lib/components/PhotosView.svelte";
  import type { Photo, Object, Tag, Metadata } from "../../../types";
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
  let sidebarObject = $state<Object | null>(null);
  let sidebarObjectTitle = $state("");
  let sidebarObjectCaption = $state("");
  let sidebarObjectTags = $state("");

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

<div class="pl-2 py-2 overflow-none">
  <h1 class="">photos</h1>
  <div class="w-full flex flex-row">
    <div class="overflow-none px-4 mt-4 w-full">
      <h3 class="text-xl font-semibold">public photos</h3>
      <div class="overflow-y-auto max-h-[40vh]">
        <PhotosView {photos} />
      </div>
      <h3 class="text-xl font-semibold mt-8">private photos</h3>
      <div class="overflow-y-auto max-h-[40vh]">
        <div
          class="flex flex-wrap w-full h-[90%] justify-center gap-4 overflow-y-auto"
        >
          {#each privatePhotos as obj}
            <div class="relative group max-w-[400px] max-h-[300px]">
              <img
                src={obj.public_url}
                alt={obj.name}
                class="w-full h-full object-cover rounded-lg"
                onclick={() => {
                  sidebarObject = obj;
                }}
              />
            </div>
          {/each}
        </div>
      </div>
    </div>
    {#if sidebarObject}
      <div
        class="w-1/3 flex flex-col items-center px-1 py-2 bg-amber-100 border-l border-amber-300 rounded-l-2xl"
      >
        <img
          src={sidebarObject.public_url}
          alt={sidebarObject.name}
          class="max-w-[90%] max-h-[90%] object-contain rounded-lg shadow-lg"
        />
        <form
          onsubmit={() => {
            const tags = sidebarObjectTags.split(",").map((tag) => tag.trim());
            fetchBackend("/api/v1/photos", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              credentials: "include",
              body: JSON.stringify({
                title: sidebarObjectTitle,
                photo_url: sidebarObject!.public_url,
                comment: sidebarObjectCaption,
                tags: tags,
              }),
            }).then(async (response) => {
              if (response.ok) {
                // remove the object from privatePhotos
                if (privatePhotos) {
                  privatePhotos = privatePhotos.filter(
                    (obj) => obj.public_url !== sidebarObject?.public_url
                  );
                }
                const newPhotoID: {
                  photo_id: number;
                  metadata: Metadata;
                } = await response.json();
                const newPhoto: Photo = {
                  id: newPhotoID.photo_id,
                  title: sidebarObjectTitle,
                  photo_url: sidebarObject!.public_url,
                  comment: sidebarObjectCaption,
                  metadata: newPhotoID.metadata,
                  tags: tags.map(
                    (tagTitle) => ({ title: tagTitle, comment: "" }) as Tag
                  ),
                };
                if (photos == undefined) {
                  photos = [newPhoto];
                } else {
                  photos = [...photos, newPhoto];
                }
                sidebarObject = null;
                sidebarObjectTitle = "";
                sidebarObjectCaption = "";
                sidebarObjectTags = "";
              } else {
                const errorText = await response.text();
                alert("error publishing photo: " + errorText);
              }
            });
          }}
          class="flex flex-col items-center w-full"
        >
          <input
            type="text"
            class="w-4/5 mt-4 px-2 py-1 rounded-lg border border-amber-900"
            placeholder={sidebarObject.name}
            bind:value={sidebarObjectTitle}
          />
          <input
            type="text"
            placeholder="caption"
            class="w-4/5 mt-4 px-2 py-1 rounded-lg border border-amber-900"
            bind:value={sidebarObjectCaption}
          />
          <input
            type="text"
            placeholder="tags (comma separated)"
            class="w-4/5 mt-4 px-2 py-1 rounded-lg border border-amber-900"
            bind:value={sidebarObjectTags}
          />
          <button
            class="bg-amber-900 text-amber-100 px-4 py-1 rounded-lg mt-4"
            type="submit">publish</button
          >
        </form>
      </div>
    {/if}
  </div>
</div>
