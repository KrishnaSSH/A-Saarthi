<script>
  import { lang, languages, t } from './lib/i18n.js';
  import Mascot from './lib/Mascot.svelte';

  let step = $state(1);
  let form = $state({
    name: '', dob: '', mobile: '', pincode: '',
    address: { line: '', village: '', district: '', state: '', pincode: '', lat: 0, lon: 0, detected: false }
  });
  let detecting = $state(false);
  let locError = $state('');
  let result = $state(null);
  let reasons = $state([]);

  const mood = $derived(detecting ? 'think' : result ? (result.ok ? 'happy' : 'oops') : 'hi');
  const speech = $derived(
    detecting ? $t('say_det') : result ? (result.ok ? $t('say_ok') : $t('say_bad')) : $t('say_s' + step)
  );

  async function detect() {
    locError = '';
    detecting = true;
    try {
      const pos = await new Promise((res, rej) =>
        navigator.geolocation.getCurrentPosition(res, rej, { enableHighAccuracy: true, timeout: 15000 }));
      const r = await fetch(`/api/location?lat=${pos.coords.latitude}&lon=${pos.coords.longitude}&lang=${$lang}`);
      const body = await r.json();
      if (!r.ok) { locError = body.why; return; }
      form.address = body;
      if (body.pincode) form.pincode = body.pincode;
    } catch {
      locError = $t('denied');
    } finally {
      detecting = false;
    }
  }

  function go(n) { step = n; result = null; }

  async function submit() {
    const r = await fetch(`/api/check?lang=${$lang}`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(form)
    });
    result = await r.json();
  }

  async function loadReasons(l) {
    try { reasons = await (await fetch(`/api/reasons?lang=${l}`)).json(); } catch { reasons = []; }
  }
  $effect(() => { loadReasons($lang); });
</script>

<a class="skip" href="#main">{$t('skip')}</a>

<div class="page">
  <header>
    <span class="logo">{$t('title')}</span>
    <label class="lang">
      <span class="sr">Language</span>
      <select bind:value={$lang}>
        {#each Object.entries(languages) as [code, name]}<option value={code}>{name}</option>{/each}
      </select>
    </label>
  </header>

  <section class="hero">
    <div class="mascot"><Mascot {mood} /></div>
    <p class="bubble" aria-live="polite">{speech}</p>
  </section>

  <main id="main">
    <div class="card">
      <ol class="dots" aria-label="Progress">
        {#each [1, 2, 3] as n}
          <li class:on={step === n} class:done={step > n}>
            <span>{step > n ? '✓' : n}</span><em>{$t('step' + n)}</em>
          </li>
        {/each}
      </ol>

      <form onsubmit={(e) => e.preventDefault()}>
        {#if step === 1}
          <label>{$t('name')}<input bind:value={form.name} autocomplete="name" /></label>
          <label>{$t('dob')}<input bind:value={form.dob} inputmode="numeric" placeholder="01-01-1990" /></label>
          <label>{$t('mobile')}<input bind:value={form.mobile} inputmode="tel" maxlength="10" autocomplete="tel-national" /></label>
          <div class="row end"><button class="primary" onclick={() => go(2)}>{$t('next')} →</button></div>
        {:else if step === 2}
          <button class="detect" onclick={detect} disabled={detecting}>📍 {detecting ? $t('detecting') : $t('detect')}</button>
          <p class="hint">{$t('detectHint')}</p>
          {#if locError}<p class="err" role="alert">{locError}</p>{/if}
          <div class="two">
            <label>{$t('village')}<input bind:value={form.address.village} /></label>
            <label>{$t('district')}<input bind:value={form.address.district} /></label>
            <label>{$t('state')}<input bind:value={form.address.state} /></label>
            <label>{$t('pincode')}<input bind:value={form.pincode} inputmode="numeric" maxlength="6" /></label>
          </div>
          <div class="row">
            <button onclick={() => go(1)}>← {$t('back')}</button>
            <button class="primary" onclick={() => go(3)}>{$t('next')} →</button>
          </div>
        {:else}
          <dl>
            <dt>{$t('name')}</dt><dd>{form.name || '—'}</dd>
            <dt>{$t('dob')}</dt><dd>{form.dob || '—'}</dd>
            <dt>{$t('mobile')}</dt><dd>{form.mobile || '—'}</dd>
            <dt>{$t('step2')}</dt><dd>{[form.address.village, form.address.district, form.address.state, form.pincode].filter(Boolean).join(', ') || '—'}</dd>
          </dl>
          <div class="row">
            <button onclick={() => go(2)}>← {$t('back')}</button>
            <button class="primary" onclick={submit}>{$t('submit')}</button>
          </div>
          {#if result}
            <div role="status" class="result">
              {#if result.ok}
                <p class="ok">🎉 {$t('allGood')}</p>
              {:else}
                <h3>{$t('problems')}</h3>
                {#each result.issues as i}
                  <article>
                    <h4>{i.title} <code>{i.code}</code></h4>
                    <p><b>{$t('why')}</b> {i.why}</p>
                    <p><b>{$t('fix')}</b> {i.fix}</p>
                  </article>
                {/each}
              {/if}
            </div>
          {/if}
        {/if}
      </form>
    </div>

    <section class="ref">
      <h2>{$t('reasons')}</h2>
      <p class="sub">{$t('reasonsSub')}</p>
      <div class="cards">
        {#each reasons as r}
          <article>
            <code>{r.code}</code>
            <h3>{r.title}</h3>
            <p>{r.why}</p>
            <p class="fix">✔ {r.fix}</p>
          </article>
        {/each}
      </div>
    </section>
  </main>
</div>

<style>
  .skip { position: absolute; left: -999px; }
  .skip:focus { left: 8px; top: 8px; background: #fff; padding: 8px; z-index: 9; }
  .sr { position: absolute; left: -999px; }
  .page { max-width: 880px; margin: 0 auto; padding: 0 20px 80px; }

  header { display: flex; justify-content: space-between; align-items: center; padding: 20px 0; }
  .logo { font-size: 1.6rem; font-weight: 900; color: var(--violet); letter-spacing: -.01em; }
  select { border: 2.5px solid var(--ink); background: #fff; padding: 6px 12px; border-radius: 999px; font-weight: 700; }

  .hero { display: flex; align-items: center; gap: 20px; margin: 8px 0 28px; }
  .mascot { width: clamp(96px, 22vw, 160px); flex: none; }
  .bubble {
    position: relative; background: #fff; border: 2.5px solid var(--ink); border-radius: 22px;
    padding: 16px 20px; font-size: clamp(1.1rem, 3.2vw, 1.5rem); font-weight: 800; line-height: 1.3;
    box-shadow: 5px 5px 0 var(--sun);
  }
  .bubble::before {
    content: ''; position: absolute; left: -12px; top: 50%; width: 20px; height: 20px; background: #fff;
    border-left: 2.5px solid var(--ink); border-bottom: 2.5px solid var(--ink); transform: translateY(-50%) rotate(45deg);
  }

  .card, .ref article { background: #fff; border: 2.5px solid var(--ink); border-radius: 24px; box-shadow: 6px 6px 0 var(--violet); }
  .card { padding: 24px; }

  .dots { list-style: none; padding: 0; display: flex; gap: 8px; margin-bottom: 24px; }
  .dots li { flex: 1; display: flex; align-items: center; gap: 8px; color: var(--mute); font-weight: 700; font-size: .8rem; }
  .dots span { width: 30px; height: 30px; border-radius: 50%; border: 2.5px solid currentColor; display: grid; place-items: center; flex: none; }
  .dots li.on { color: var(--ink); }
  .dots li.on span { background: var(--sun); border-color: var(--ink); }
  .dots li.done { color: var(--leaf); }
  .dots li.done span { background: var(--leaf); border-color: var(--leaf); color: #fff; }
  .dots em { font-style: normal; }

  form { display: grid; gap: 20px; }
  label { display: grid; gap: 6px; font-weight: 800; font-size: .9rem; }
  input {
    border: 2.5px solid var(--ink); border-radius: 14px; padding: 12px 14px; background: var(--bg); font-size: 1.1rem; font-weight: 600;
  }
  input:focus { outline: none; border-color: var(--violet); box-shadow: 0 0 0 4px var(--lav); background: #fff; }
  .two { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

  button {
    border: 2.5px solid var(--ink); background: #fff; padding: 12px 22px; font-weight: 800; cursor: pointer;
    border-radius: 999px; min-height: 48px; box-shadow: 3px 3px 0 var(--ink); transition: transform .08s, box-shadow .08s;
  }
  button:hover:not(:disabled) { transform: translate(-1px, -1px); box-shadow: 5px 5px 0 var(--ink); }
  button:active:not(:disabled) { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
  button.primary { background: var(--violet); color: #fff; }
  button.detect { background: var(--sun); width: 100%; font-size: 1.1rem; }
  button:disabled { opacity: .6; cursor: wait; }
  .row { display: flex; justify-content: space-between; gap: 12px; flex-wrap: wrap; }
  .row.end { justify-content: flex-end; }
  .hint { color: var(--mute); margin-top: -8px; font-size: .9rem; }
  .err { color: #c0301f; font-weight: 700; }

  dl { display: grid; grid-template-columns: 1fr 2fr; gap: 10px 20px; }
  dt { color: var(--mute); font-size: .85rem; } dd { font-weight: 800; }

  .result { display: grid; gap: 14px; border-top: 3px dashed var(--lav); padding-top: 18px; }
  .ok { font-size: 1.5rem; font-weight: 900; color: var(--leaf); }
  .result article { background: #fff1ef; border: 2.5px solid var(--coral); border-radius: 16px; padding: 14px 16px; display: grid; gap: 4px; }
  h4 { font-size: 1.05rem; }
  code { font-family: ui-monospace, monospace; font-size: .7rem; background: var(--lav); color: var(--violet); padding: 2px 8px; border-radius: 999px; margin-left: 6px; }

  .ref { margin-top: 64px; }
  h2 { font-size: clamp(1.6rem, 5vw, 2.2rem); font-weight: 900; letter-spacing: -.01em; }
  .sub { color: var(--mute); margin: 4px 0 24px; }
  .cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); gap: 22px; }
  .cards article { padding: 18px; display: grid; gap: 8px; align-content: start; box-shadow: 4px 4px 0 var(--ink); }
  .cards article:nth-child(3n+1) { box-shadow: 4px 4px 0 var(--sun); }
  .cards article:nth-child(3n+2) { box-shadow: 4px 4px 0 var(--coral); }
  .cards article code { margin: 0; justify-self: start; }
  .cards h3 { font-size: 1.05rem; line-height: 1.25; }
  .cards p { font-size: .9rem; }
  .fix { color: var(--leaf); font-weight: 700; }

  @media (max-width: 600px) {
    .hero { gap: 16px; }
    .dots em { display: none; }
    .two { grid-template-columns: 1fr; }
    dl { grid-template-columns: 1fr; gap: 2px; } dd { margin-bottom: 10px; }
  }
</style>
