<script>
  // Saarthi: the hand-drawn circle from the original doodle, given a face.
  let { mood = 'hi' } = $props();
</script>

<svg class="m {mood}" viewBox="-6 -6 112 106" role="img" aria-label="Saarthi">
  <g class="body">
    <path
      fill="#efeaff" stroke="#6741d9" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"
      d="M67.72 3.05 C75.38 4.97, 83.41 10.46, 89 16.46 C94.58 22.46, 99.76 31.34, 101.23 39.04 C102.7 46.75, 100.77 55.32, 97.84 62.68 C94.91 70.04, 89.92 78.28, 83.65 83.19 C77.39 88.1, 68.64 90.78, 60.25 92.12 C51.87 93.46, 41.3 93.7, 33.35 91.2 C25.4 88.71, 18.12 83.26, 12.56 77.14 C7 71.02, 1.72 62.22, -0.03 54.48 C-1.77 46.74, -0.82 37.85, 2.09 30.71 C5 23.57, 10.75 16.47, 17.45 11.65 C24.14 6.83, 33.03 3.1, 42.27 1.78 C51.51 0.46, 66.98 2.8, 72.89 3.73 C78.81 4.66, 78.28 6.67, 77.75 7.38"
    />
    <path fill="none" stroke="#6741d9" stroke-width="1.5" stroke-linecap="round" opacity=".5"
      d="M60.71 0.44 C68.67 1.51, 77.57 7.23, 83.75 12.45 C89.93 17.67, 95.16 24.68, 97.81 31.76 C100.47 38.84, 101.29 47.12, 99.67 54.92 C98.06 62.73, 93.54 72.76, 88.14 78.6 C82.73 84.45, 75.1 87.32, 67.24 89.98 C59.37 92.63, 49.21 95.9, 40.96 94.53 C32.7 93.16, 23.87 86.93, 17.71 81.74" />

    <circle cx="24" cy="58" r="6" fill="#ff6b5b" opacity=".35" />
    <circle cx="78" cy="58" r="6" fill="#ff6b5b" opacity=".35" />

    <g class="eyes" stroke="#2a1f4d" stroke-width="3.5" stroke-linecap="round" fill="#2a1f4d">
      {#if mood === 'happy'}
        <path d="M30 46 Q36 36 42 46" fill="none" /><path d="M60 46 Q66 36 72 46" fill="none" />
      {:else if mood === 'think'}
        <circle class="look" cx="37" cy="43" r="4.5" stroke="none" /><circle class="look" cx="67" cy="43" r="4.5" stroke="none" />
      {:else if mood === 'oops'}
        <circle cx="36" cy="46" r="4.5" stroke="none" /><circle cx="66" cy="46" r="4.5" stroke="none" />
        <path d="M28 34 L43 38 M74 34 L59 38" fill="none" />
      {:else}
        <g class="blink"><circle cx="36" cy="44" r="4.5" stroke="none" /><circle cx="66" cy="44" r="4.5" stroke="none" /></g>
      {/if}
    </g>

    <g fill="none" stroke="#2a1f4d" stroke-width="3.5" stroke-linecap="round">
      {#if mood === 'happy'}
        <path d="M36 62 Q51 80 66 62 Z" fill="#ff6b5b" stroke-linejoin="round" />
      {:else if mood === 'think'}
        <ellipse cx="51" cy="68" rx="4" ry="5" />
      {:else if mood === 'oops'}
        <path d="M40 72 Q46 64 51 70 Q56 76 62 68" />
      {:else}
        <path d="M38 62 Q51 74 64 62" />
      {/if}
    </g>
  </g>
</svg>

<style>
  .m { width: 100%; height: auto; overflow: visible; }
  .body { transform-origin: 50% 90%; transform-box: fill-box; }
  .hi .body { animation: bob 3.2s ease-in-out infinite; }
  .think .body { animation: wiggle .6s ease-in-out infinite; }
  .happy .body { animation: hop .7s ease-out 2; }
  .oops .body { animation: shake .5s ease-in-out 1; }
  .blink { animation: blink 4s infinite; transform-origin: center; transform-box: fill-box; }
  @keyframes bob { 50% { transform: translateY(-4px) rotate(1.5deg); } }
  @keyframes wiggle { 25% { transform: rotate(-5deg); } 75% { transform: rotate(5deg); } }
  @keyframes hop { 40% { transform: translateY(-14px) scale(1.04); } }
  @keyframes shake { 20%, 60% { transform: translateX(-4px); } 40%, 80% { transform: translateX(4px); } }
  @keyframes blink { 94%, 100% { transform: scaleY(1); } 97% { transform: scaleY(.1); } }
</style>
