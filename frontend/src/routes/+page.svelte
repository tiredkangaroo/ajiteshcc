<script>
  import { isAdmin } from "../lib/stores/isAdmin";
</script>

<div class="w-full h-full flex flex-col justify-center align-middle p-4">
  <div class="flex flex-row justify-between items-center px-12">
    <h2><a href="/projects">Projects</a></h2>
    <div class="text-center">
      <img
        src="https://avatars.githubusercontent.com/u/81335306?v=4"
        alt="tiredkangaroo pfp"
        class="w-48"
      />
      <h1 class="mt-4">{$isAdmin ? "hi, aji!" : "Ajitesh Kumar"}</h1>
    </div>
    <h2><a href="/photography">Photography</a></h2>
  </div>
</div>
