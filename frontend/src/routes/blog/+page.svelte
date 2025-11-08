<script lang="ts">
  import type { ErrorResponse, PostHead } from "../../types";
  import { fetchBackend, formatDate } from "../../utils";

  let postHeads: PostHead[] | undefined = $state(undefined);
  $effect(() => {
    async function loadPosts() {
      const res = await fetchBackend("/api/v1/posts");
      const data: PostHead[] | ErrorResponse = await res.json();
      if (res.status !== 200) {
        console.error("fetch posts", (data as ErrorResponse).error);
        return;
      }
      if (data === null) {
        postHeads = [];
        return;
      }
      postHeads = data as PostHead[];
    }
    loadPosts();
  });
</script>

<div class="w-full h-full flex flex-col items-center">
  <div class="w-full relative flex items-center justify-center mt-3 mb-4">
    <button
      onclick={() => {
        window.location.replace("/");
      }}
      class="absolute left-4 px-4 py-1 bg-amber-400">back</button
    >
    <h1 class="text-xl font-semibold">blog 📝</h1>
  </div>
  {#if postHeads == undefined}
    <div class="w-full h-full justify-center flex items-center align-middle">
      <p>loading blog</p>
    </div>
  {:else if postHeads.length === 0}
    <h3
      class="w-full h-full justify-center flex items-center align-middle font-light text-xl"
    >
      nothing here right now
    </h3>
  {:else}
    <div class="w-[50%] min-w-[300px] justify-center">
      {#each postHeads as postHead}
        <div
          class="w-full border-t-[0.5px] border-t-black border-b-[0.5px] border-b-black flex flex-row"
        >
          <div class="flex flex-col flex-1">
            <a href={`/blog/${postHead.slug}`} class="text-2xl font-semibold"
              >{postHead.title}</a
            >
            <span class="text-sm font-light text-gray-500"
              >{formatDate(new Date(postHead.created_at))}</span
            >
            <p class="mt-2">{postHead.comment}</p>
          </div>
          <div class="flex-1 flex flex-col items-end py-1 justify-between">
            <span
              class="text-sm font-light"
              style={postHead.published ? "color: #02a317" : "color: #b30015"}
              >{postHead.published ? "published" : "unpublished"}</span
            >
            <div class="flex flex-wrap justify-end">
              {#each postHead.tags as tag}
                <button
                  class="ml-2 px-4 h-fit bg-amber-300 text-sm font-light"
                  title={tag.comment}
                  onclick={() => {
                    window.location.replace(`/tags/${tag.title}`);
                  }}>{tag.title}</button
                >
              {/each}
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
