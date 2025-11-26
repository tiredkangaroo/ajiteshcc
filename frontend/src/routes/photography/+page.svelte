<script lang="ts">
  import PhotosView from "$lib/components/PhotosView.svelte";
  import type { Photo, Tag, Metadata } from "../../types";
  import { fetchBackend } from "../../utils";

  let photos: Photo[] | null = $state(null);
  let photoCols: Photo[][] = $derived.by(() => {
    if (photos === null) {
      return [];
    }
    let cols = [];
    for (let i = 0; i < photos.length; i += 3) {
      cols.push(photos.slice(i, i + 3));
    }
    return cols;
  });

  let lastClickedImageID: number | undefined = $state(undefined);
  let lcImageMetadata: Metadata | null = $derived.by(() => {
    if (lastClickedImageID === undefined || photos === null) {
      return null;
    }
    if (photos[lastClickedImageID] && photos[lastClickedImageID].metadata) {
      return photos[lastClickedImageID].metadata!;
    }
    return null;
  });

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
  function metadata(md: Metadata): string {
    return (
      metadataLine1(md) + "\n" + metadataLine2(md) + "\n" + metadataLine3(md)
    );
  }
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
  <div class="h-screen">
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
