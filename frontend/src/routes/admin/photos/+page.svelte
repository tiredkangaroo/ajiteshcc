<script lang="ts">
  import type { Photo, Object, Tag, Metadata } from "../../../types";
  import { fetchBackend } from "../../../utils";

  let photos = $state<Array<Photo> | undefined>(undefined);
  let objects = $state<Array<Object> | undefined>(undefined);
  let privatePhotos: Array<Object> | undefined = $derived.by(() => {
    if (objects == undefined || photos == undefined) {
      return undefined;
    }
    const photoUrls = new Set(photos.map((photo) => photo.photo_url));
    return objects.filter((obj) => !photoUrls.has(obj.public_url));
  });

  let sidebarObject = $state<Object | null>(null);
  let sidebarObjectTitle = $state("");
  let sidebarObjectCaption = $state("");
  let sidebarObjectTags = $state("");

  let sidebarPhoto = $state<Photo | null>(null);

  let showingSidebarNewPhoto = $state(false);
  let sidebarNewPhotoFile: File | null = $state(null);
  let sidebarNewPhotoTitle = $state("");
  let sidebarNewPhotoCaption = $state("");
  let sidebarNewPhotoTags = $state("");
  let sidebarNewPhotoPublic = $state(true);

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

<div
  class="min-h-screen bg-gradient-to-br from-amber-50 via-orange-50 to-rose-50"
>
  <div class="max-w-[1800px] mx-auto px-6 py-8">
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-2">
        photography admin dashboard
      </h1>
    </div>

    <div class="justify-end flex mb-2">
      <button
        class="bg-amber-800 text-amber-100 px-6 py-1 rounded-xl"
        onclick={() => {
          showingSidebarNewPhoto = true;
        }}>upload</button
      >
    </div>
    <div class="flex gap-6">
      <div class="flex-1 space-y-8">
        <div class="bg-white rounded-2xl shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-2xl font-bold text-gray-800">public</h2>
              <p class="text-sm text-gray-500 mt-1">
                {photos?.length || 0} photos published
              </p>
            </div>
          </div>

          <div
            class="flex flex-wrap w-full max-h-[50vh] justify-center gap-4 overflow-y-auto"
          >
            {#each photos as photo}
              <div
                class="relative group max-w-[400px] max-h-[300px]"
                onclick={() => {
                  if (sidebarObject != null) {
                    sidebarObject = null;
                  }
                  sidebarPhoto = photo;
                }}
                role="button"
                tabindex="0"
                onkeypress={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    if (sidebarObject != null) {
                      sidebarObject = null;
                    }
                    sidebarPhoto = photo;
                  }
                }}
              >
                <img
                  src={photo.photo_url}
                  alt={photo.title}
                  class="w-full h-full object-cover rounded-lg shadow-md group-hover:shadow-xl transition-shadow"
                />
                <div
                  class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-xl flex items-end p-3"
                >
                  <p class="text-white text-sm font-medium truncate w-full">
                    {photo.title}
                  </p>
                </div>
              </div>
            {/each}
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-lg p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-2xl font-bold text-gray-800">private photos</h2>
              <p class="text-sm text-gray-500 mt-1">
                {privatePhotos?.length || 0} photos not published
              </p>
            </div>
            <div
              class="w-12 h-12 bg-amber-100 rounded-full flex items-center justify-center"
            >
              <svg
                class="w-6 h-6 text-amber-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                ></path>
              </svg>
            </div>
          </div>

          <div
            class="flex flex-wrap w-full max-h-[50vh] justify-center gap-4 overflow-y-auto"
          >
            {#each privatePhotos as obj}
              <div
                class="relative group max-w-[400px] max-h-[300px]"
                onclick={() => {
                  if (sidebarPhoto != null) {
                    sidebarPhoto = null;
                  }
                  sidebarObject = obj;
                }}
              >
                <img
                  src={obj.public_url}
                  alt={obj.name}
                  class="w-full h-full object-cover rounded-lg shadow-md group-hover:shadow-xl transition-shadow"
                />
                <div
                  class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-xl flex items-end p-3"
                >
                  <p class="text-white text-sm font-medium truncate w-full">
                    {obj.name}
                  </p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>

      <!-- Sidebars -->
      {#if sidebarObject}
        {@render objectSidebar()}
      {/if}
      {#if sidebarPhoto}
        {@render photoSidebar()}
      {/if}
      {#if showingSidebarNewPhoto}
        {@render newPhotoSidebar()}
      {/if}
    </div>
  </div>
</div>

{#snippet objectSidebar()}
  <div
    class="w-[400px] bg-amber-100 border border-amber-300 rounded-2xl shadow-2xl p-6 sticky top-8 h-fit overflow-y-auto"
  >
    <div class="flex justify-between items-center mb-6">
      <h3 class="text-xl font-bold text-gray-800">publish photo</h3>
      <button
        class="w-8 h-8 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center transition-colors"
        onclick={() => {
          sidebarObject = null;
        }}
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          ></path>
        </svg>
      </button>
    </div>

    <div class="mb-6">
      <img
        src={sidebarObject!.public_url}
        alt={sidebarObject!.name}
        class="w-full h-48 object-contain rounded-xl"
      />
    </div>

    <form
      onsubmit={() => {
        const tags = sidebarObjectTags.split(",").map((tag) => tag.trim());
        fetchBackend("/api/v1/photos", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({
            title: sidebarObjectTitle,
            photo_url: sidebarObject!.public_url,
            comment: sidebarObjectCaption,
            tags: tags,
          }),
        }).then(async (response) => {
          if (response.ok) {
            if (privatePhotos) {
              privatePhotos = privatePhotos.filter(
                (obj) => obj.public_url !== sidebarObject?.public_url
              );
            }
            const newPhotoID: { photo_id: number; metadata: Metadata } =
              await response.json();
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
            photos = photos == undefined ? [newPhoto] : [...photos, newPhoto];
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
      class="space-y-4"
    >
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Title</label
        >
        <input
          type="text"
          class="w-full px-4 py-2.5 rounded-lg border border-gray-300 focus:ring-2 focus:ring-amber-500 focus:border-transparent transition-all outline-none"
          placeholder={sidebarObject!.name}
          bind:value={sidebarObjectTitle}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2"
          >Caption</label
        >
        <textarea
          class="w-full px-4 py-2.5 rounded-lg border border-gray-300 focus:ring-2 focus:ring-amber-500 focus:border-transparent transition-all outline-none resize-none"
          rows="3"
          placeholder="Add a caption..."
          bind:value={sidebarObjectCaption}
        ></textarea>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Tags</label>
        <input
          type="text"
          placeholder="nature, landscape, sunset"
          class="w-full px-4 py-2.5 rounded-lg border border-gray-300 focus:ring-2 focus:ring-amber-500 focus:border-transparent transition-all outline-none"
          bind:value={sidebarObjectTags}
        />
        <p class="text-xs text-gray-500 mt-1">Separate tags with commas</p>
      </div>

      <button
        class="w-full bg-amber-900 text-amber-200 px-6 py-3 rounded-lg shadow-md hover:shadow-lg transition-all"
        type="submit"
      >
        publish
      </button>
    </form>
  </div>
{/snippet}

{#snippet photoSidebar()}
  <div
    class="w-[400px] bg-amber-100 border border-amber-300 rounded-2xl shadow-2xl p-6 sticky top-8 h-fit overflow-y-auto"
  >
    <div class="flex justify-between items-center mb-6">
      <h3 class="text-xl font-bold text-gray-800">photo details</h3>
      <button
        class="w-8 h-8 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center transition-colors"
        onclick={() => {
          sidebarPhoto = null;
        }}
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          ></path>
        </svg>
      </button>
    </div>

    <div class="mb-6">
      <img
        src={sidebarPhoto!.photo_url}
        alt={sidebarPhoto!.title}
        class="w-full h-48 object-contain rounded-xl"
      />
    </div>

    <div class="space-y-4">
      <div>
        <h4 class="text-lg font-semibold text-gray-800">
          {sidebarPhoto!.title}
        </h4>
      </div>

      {#if sidebarPhoto!.comment}
        <div>
          <p class="text-gray-700">{sidebarPhoto!.comment}</p>
        </div>
      {/if}

      <div>
        <label class="block text-sm font-medium text-gray-600 mb-2">Tags</label>
        <div class="flex flex-wrap gap-2">
          {#each sidebarPhoto!.tags as tag}
            <span
              class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gradient-to-r from-amber-100 to-orange-100 text-amber-800 border border-amber-200"
            >
              {tag.title}
            </span>
          {/each}
        </div>
      </div>

      <button
        class="w-full bg-amber-900 text-amber-200 px-6 py-3 rounded-lg shadow-md hover:shadow-lg transition-all"
        onclick={() => {
          fetchBackend("/api/v1/photos/" + sidebarPhoto!.id, {
            method: "DELETE",
            credentials: "include",
          }).then(async (response) => {
            if (response.ok) {
              photos = photos?.filter((photo) => photo.id !== sidebarPhoto?.id);
              sidebarPhoto = null;
            } else {
              const errorText = await response.text();
              alert("error unpublishing photo: " + errorText);
            }
          });
        }}
      >
        unpublish
      </button>
    </div>
  </div>
{/snippet}

{#snippet newPhotoSidebar()}
  <div
    class="w-[400px] bg-amber-100 border border-amber-300 rounded-2xl shadow-2xl p-6 sticky top-8 h-fit overflow-y-auto"
  >
    <div class="flex justify-between items-center mb-6">
      <h3 class="text-xl font-bold text-gray-800">add new photo</h3>
      <button
        class="w-8 h-8 rounded-full bg-red-500 hover:bg-red-600 text-white flex items-center justify-center transition-colors"
        onclick={() => {
          showingSidebarNewPhoto = false;
        }}
        title="close new photo sidebar"
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          ></path>
        </svg>
      </button>
    </div>

    <!-- Drag and Drop Area -->
    {#if sidebarNewPhotoFile}
      <div class="mb-6 flex flex-col justify-center items-center gap-2">
        <img
          src={URL.createObjectURL(sidebarNewPhotoFile!)}
          alt="preview"
          class="w-full h-48 object-contain rounded-xl"
        />
        <button
          onclick={() => {
            sidebarNewPhotoFile = null;
          }}
          class="bg-amber-800 text-amber-100 px-3 py-1 rounded-xl"
        >
          remove file
        </button>
      </div>
    {:else}
      <div
        class="border-2 border-dashed border-gray-300 rounded-lg p-8 mb-6 text-center hover:border-blue-400 transition-colors cursor-pointer"
        ondragover={(event) => {
          event.preventDefault();
          event.currentTarget.classList.add("border-blue-500", "bg-blue-50");
        }}
        ondragleave={(event) => {
          event.currentTarget.classList.remove("border-blue-500", "bg-blue-50");
        }}
        ondrop={(event) => {
          event.preventDefault();
          event.currentTarget.classList.remove("border-blue-500", "bg-blue-50");
          const file = event.dataTransfer.files[0];
          if (file && file.type.startsWith("image/")) {
            sidebarNewPhotoFile = file;
          }
        }}
        onclick={() => {
          document.getElementById("fileInput")!.click();
        }}
      >
        <svg
          class="w-12 h-12 mx-auto mb-3 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
          ></path>
        </svg>
        <p class="text-gray-600 mb-2">Drag and drop your photo here</p>
        <p class="text-sm text-gray-500">or click to browse</p>
        <input
          id="fileInput"
          type="file"
          accept="image/*"
          class="hidden"
          onchange={(e) => {
            const file = e.target!.files[0];
            if (file) {
              sidebarNewPhotoFile = file;
            }
          }}
        />
      </div>
    {/if}

    <div class="space-y-4">
      <div class="flex flex-row items-center gap-2">
        <p>public?</p>
        <div class="relative inline-block w-11 h-5">
          <input
            id="switch-component"
            type="checkbox"
            class="peer appearance-none w-11 h-5 bg-slate-400 rounded-full checked:bg-amber-800 cursor-pointer transition-colors duration-300"
            bind:checked={sidebarNewPhotoPublic}
          />
          <label
            for="switch-component"
            class="absolute top-0 left-0 w-5 h-5 bg-white rounded-full border border-slate-300 shadow-sm transition-transform duration-300 peer-checked:translate-x-6 peer-checked:border-slate-800 cursor-pointer"
          >
          </label>
        </div>
      </div>
      {#if sidebarNewPhotoPublic}
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Title
          </label>
          <input
            type="text"
            bind:value={sidebarNewPhotoTitle}
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter photo title"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Caption
          </label>
          <textarea
            bind:value={sidebarNewPhotoCaption}
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
            placeholder="Enter photo caption"
          ></textarea>
        </div>

        <!-- Tags Field -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Tags
          </label>
          <input
            type="text"
            bind:value={sidebarNewPhotoTags}
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g. nature, sunset, landscape"
          />
        </div>
      {/if}

      <!-- Submit Button -->
      <button
        class="w-full bg-amber-900 text-amber-200 px-6 py-3 rounded-lg shadow-md hover:shadow-lg transition-all"
        onclick={() => {
          // need to do two things: upload the file /api/v1/objects then publish the photo /api/v1/photos
          if (sidebarNewPhotoFile == null) {
            alert("please select a photo file to upload");
            return;
          }
          const formData = new FormData();
          formData.append("file", sidebarNewPhotoFile!);
          formData.append("name", sidebarNewPhotoFile!.name);
          fetchBackend("/api/v1/objects", {
            method: "POST",
            body: formData,
            credentials: "include",
          }).then(async (response) => {
            if (response.ok) {
              const newObject = await response.json();
              if (sidebarNewPhotoPublic) {
                // publish photo
                const tags = sidebarNewPhotoTags
                  .split(",")
                  .map((tag) => tag.trim());
                const photoResponse = await fetchBackend("/api/v1/photos", {
                  method: "POST",
                  headers: { "Content-Type": "application/json" },
                  credentials: "include",
                  body: JSON.stringify({
                    title: sidebarNewPhotoTitle,
                    photo_url: newObject.public_url,
                    comment: sidebarNewPhotoCaption,
                    tags: tags,
                  }),
                });
                if (photoResponse.ok) {
                  const newPhotoID: {
                    photo_id: number;
                    metadata: Metadata;
                  } = await photoResponse.json();
                  const newPhoto: Photo = {
                    id: newPhotoID.photo_id,
                    title: sidebarNewPhotoTitle,
                    photo_url: newObject.public_url,
                    comment: sidebarNewPhotoCaption,
                    metadata: newPhotoID.metadata,
                    tags: tags.map(
                      (tagTitle) => ({ title: tagTitle, comment: "" }) as Tag
                    ),
                  };
                  photos =
                    photos == undefined ? [newPhoto] : [...photos, newPhoto];
                  showingSidebarNewPhoto = false;
                  sidebarNewPhotoFile = null;
                  sidebarNewPhotoTitle = "";
                  sidebarNewPhotoCaption = "";
                  sidebarNewPhotoTags = "";
                } else {
                  const errorJSON = await photoResponse.json();
                  alert("error publishing photo: " + errorJSON.error);
                }
              } else {
                // add private photo and close the sidebar
                privatePhotos =
                  privatePhotos == undefined
                    ? [newObject]
                    : [...privatePhotos, newObject];
                showingSidebarNewPhoto = false;
                sidebarNewPhotoFile = null;
                sidebarNewPhotoTitle = "";
                sidebarNewPhotoCaption = "";
                sidebarNewPhotoTags = "";
              }
            } else {
              const errorJSON = await response.json();
              alert("error uploading photo file: " + errorJSON.error);
            }
          });
        }}
      >
        {#if sidebarNewPhotoPublic}
          publish
        {:else}
          upload
        {/if}
      </button>
    </div>
  </div>
{/snippet}
