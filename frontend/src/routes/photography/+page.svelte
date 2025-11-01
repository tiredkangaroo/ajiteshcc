<script lang="ts">
  import type { Photo, Tag, Metadata } from "../../types";
  import { fetchBackend } from "../../utils";

  let photos: Photo[] | null = $state(null);
  $effect(() => {
    async function loadPhotos() {
      const res = await fetchBackend("/api/v1/photos");
      if (res.ok) {
        const data = await res.json();
        console.log(data);
        photos = data;
      } else {
        console.error("failed to load photos", await res.text());
      }
    }
    loadPhotos();
  });

  const isValid = (v: string): boolean =>
    v !== null &&
    v !== undefined &&
    v.trim().length > 0 &&
    v !== "0.000000" &&
    v !== "0001-01-01 00:00:00 +0000 UTC" &&
    v !== "0";
  const onlyValid = (v: string): string => (isValid(v) ? v : "");
  // metadata:
  // date @ lat long altitude
  // camera make model lens make model focal length
  // iso, ss, aperture
  function metadataLine1(md: Metadata): string {
    const parts: string[] = [];
    if (isValid(md.createdat)) {
      const d = new Date(md.createdat);
      parts.push(
        [
          "Sunday",
          "Monday",
          "Tuesday",
          "Wednesday",
          "Thursday",
          "Friday",
          "Saturday",
        ][d.getDay()] +
          ", " +
          d.toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour12: true,
            hour: "2-digit",
            minute: "2-digit",
            timeZoneName: "short",
          })
      );
    }
    if (isValid(md.latitude) && isValid(md.longitude)) {
      parts.push(
        `@ ${onlyValid(md.latitude)}, ${onlyValid(md.longitude)}${
          isValid(md.altitude) ? `, ${md.altitude}` : ""
        }`
      );
    }
    return parts.join(" ");
  }
  function metadataLine2(md: Metadata): string {
    const parts: string[] = [];
    if (isValid(md.cameramake) || isValid(md.cameramodel)) {
      parts.push(
        `Camera: ${onlyValid(md.cameramake)} ${onlyValid(md.cameramodel)}`.trim()
      );
    }
    if (isValid(md.lensmake) || isValid(md.lensmodel)) {
      parts.push(
        `Lens: ${onlyValid(md.lensmake)} ${onlyValid(md.lensmodel)} ${onlyValid(parseFloat(md.focallength).toFixed(0) + "mm")}`.trim()
      );
    }
    return parts.join(" | ");
  }
  function metadataLine3(md: Metadata): string {
    const parts: string[] = [];
    if (isValid(md.iso)) {
      parts.push(`ISO ${md.iso}`);
    }
    if (isValid(md.shutterspeed)) {
      parts.push(`Shutter ${md.shutterspeed}`);
    }
    if (isValid(md.aperture)) {
      parts.push(`𝑓/${parseFloat(md.aperture).toFixed(0)}`);
    }
    return parts.join(", ");
  }
</script>

{#if photos === null}
  <div class="w-full h-full justify-center items-center flex">
    Loading photos...
  </div>
{:else}
  <div class="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
    {#each photos as photo}
      <div class="m-4 rounded-lg overflow-hidden bg-gray-200">
        <img
          src={photo.photo_url}
          alt={photo.title || "Photo"}
          class="w-full"
        />
        {#if photo.title || photo.comment}
          <div class="p-4">
            <div class="flex flex-row gap-4">
              {#if photo.title}
                <h3 class="text-lg font-semibold mb-2">{photo.title}</h3>
              {/if}
              <div>
                {#each photo.tags as tag}
                  <span
                    class="inline-block bg-gray-400 rounded-full px-3 py-1 text-sm font-semibold text-gray-700 mr-2 mb-2"
                    >{tag.title}</span
                  >
                {/each}
              </div>
            </div>
            {#if photo.comment}
              <p class="text-gray-700 text-base">{photo.comment}</p>
            {/if}
          </div>
          {#if photo.metadata}
            <div class="mt-3 px-4 py-2">
              <p>{metadataLine1(photo.metadata)}</p>
              <p>{metadataLine2(photo.metadata)}</p>
              <p>{metadataLine3(photo.metadata)}</p>
            </div>
          {/if}
        {/if}
      </div>
    {/each}
  </div>
{/if}
