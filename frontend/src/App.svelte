<script>
  import { lang, languages, t } from './lib/i18n.js';

  let step = $state(1);
  let form = $state({
    name: '', dob: '', mobile: '', pincode: '',
    address: { line: '', village: '', district: '', state: '', pincode: '', lat: 0, lon: 0, detected: false }
  });
  let detecting = $state(false);
  let locError = $state('');
  let result = $state(null);
  let reasons = $state([]);

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

<header>
  <div class="bar">
    <h1>{$t('title')}</h1>
    <label class="lang">
      <span class="sr">Language</span>
      <select bind:value={$lang}>
        {#each Object.entries(languages) as [code, name]}<option value={code}>{name}</option>{/each}
      </select>
    </label>
  </div>
  <p class="sub">{$t('sub')}</p>
</header>

<main id="main">
  <section class="grid">
    <ol class="steps" aria-label="Progress">
      {#each [1, 2, 3] as n}
        <li class:on={step === n} class:done={step > n}><span>{n}</span>{$t('step' + n)}</li>
      {/each}
    </ol>

    <form onsubmit={(e) => { e.preventDefault(); }}>
      {#if step === 1}
        <label>{$t('name')}<input bind:value={form.name} autocomplete="name" /></label>
        <label>{$t('dob')}<input bind:value={form.dob} inputmode="numeric" placeholder="01-01-1990" /></label>
        <label>{$t('mobile')}<input bind:value={form.mobile} inputmode="tel" maxlength="10" autocomplete="tel-national" /></label>
        <div class="row"><button class="primary" onclick={() => (step = 2)}>{$t('next')} →</button></div>
      {:else if step === 2}
        <button class="primary wide" onclick={detect} disabled={detecting}>◎ {detecting ? $t('detecting') : $t('detect')}</button>
        <p class="hint">{$t('detectHint')}</p>
        {#if locError}<p class="err" role="alert">{locError}</p>{/if}
        <label>{$t('village')}<input bind:value={form.address.village} /></label>
        <label>{$t('district')}<input bind:value={form.address.district} /></label>
        <label>{$t('state')}<input bind:value={form.address.state} /></label>
        <label>{$t('pincode')}<input bind:value={form.pincode} inputmode="numeric" maxlength="6" /></label>
        <div class="row">
          <button onclick={() => (step = 1)}>← {$t('back')}</button>
          <button class="primary" onclick={() => (step = 3)}>{$t('next')} →</button>
        </div>
      {:else}
        <dl>
          <dt>{$t('name')}</dt><dd>{form.name || '—'}</dd>
          <dt>{$t('dob')}</dt><dd>{form.dob || '—'}</dd>
          <dt>{$t('mobile')}</dt><dd>{form.mobile || '—'}</dd>
          <dt>{$t('step2')}</dt><dd>{[form.address.village, form.address.district, form.address.state, form.pincode].filter(Boolean).join(', ') || '—'}</dd>
        </dl>
        <div class="row">
          <button onclick={() => (step = 2)}>← {$t('back')}</button>
          <button class="primary" onclick={submit}>{$t('submit')}</button>
        </div>
        {#if result}
          <div role="status" class="result">
            {#if result.ok}
              <p class="ok">✓ {$t('allGood')}</p>
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
  </section>

  <section class="ref">
    <h2>{$t('reasons')}</h2>
    <p class="sub">{$t('reasonsSub')}</p>
    <table>
      <tbody>
        {#each reasons as r}
          <tr><th scope="row"><code>{r.code}</code></th><td><b>{r.title}</b><br />{r.why}<br /><span class="mute">→ {r.fix}</span></td></tr>
        {/each}
      </tbody>
    </table>
  </section>
</main>

<style>
  .skip { position: absolute; left: -999px; }
  .skip:focus { left: 8px; top: 8px; background: var(--paper); padding: 8px; z-index: 9; }
  .sr { position: absolute; left: -999px; }
  header, main { max-width: 1120px; margin: 0 auto; padding: 0 24px; }
  .bar { display: flex; justify-content: space-between; align-items: baseline; border-bottom: 4px solid var(--rule); padding: 32px 0 12px; }
  h1 { font-size: clamp(2.5rem, 8vw, 5rem); font-weight: 700; letter-spacing: -0.03em; line-height: 1; }
  .sub { color: var(--mute); margin: 12px 0 48px; }
  select { border: 2px solid var(--ink); background: var(--paper); padding: 6px 10px; border-radius: 0; }
  .grid { display: grid; grid-template-columns: 1fr 2fr; gap: 48px; margin-bottom: 96px; }
  .steps { list-style: none; padding: 0; display: grid; gap: 16px; align-content: start; }
  .steps li { display: flex; gap: 16px; align-items: center; color: var(--mute); font-weight: 600; }
  .steps li span { width: 32px; height: 32px; border: 2px solid currentColor; display: grid; place-items: center; }
  .steps li.on { color: var(--ink); }
  .steps li.on span { background: var(--accent); border-color: var(--accent); color: #fff; }
  .steps li.done span { background: var(--ink); border-color: var(--ink); color: #fff; }
  form { display: grid; gap: 24px; border-top: 4px solid var(--rule); padding-top: 24px; }
  label { display: grid; gap: 6px; font-weight: 600; }
  input { border: 0; border-bottom: 2px solid var(--ink); padding: 10px 0; font-size: 1.2rem; background: transparent; border-radius: 0; }
  input:focus { outline: none; border-bottom-color: var(--accent); border-bottom-width: 4px; }
  button { border: 2px solid var(--ink); background: var(--paper); padding: 12px 24px; font-weight: 700; cursor: pointer; border-radius: 0; min-height: 48px; }
  button.primary { background: var(--ink); color: #fff; }
  button:hover:not(:disabled) { background: var(--accent); border-color: var(--accent); color: #fff; }
  button:disabled { opacity: .5; cursor: wait; }
  .wide { width: 100%; text-align: left; }
  .row { display: flex; justify-content: space-between; gap: 16px; }
  .hint { color: var(--mute); margin-top: -12px; }
  .err { color: var(--accent); font-weight: 600; }
  dl { display: grid; grid-template-columns: 1fr 2fr; gap: 12px 24px; }
  dt { color: var(--mute); } dd { font-weight: 600; }
  .result { border-top: 4px solid var(--rule); padding-top: 16px; display: grid; gap: 16px; }
  .ok { font-size: 1.4rem; font-weight: 700; }
  article { border-left: 6px solid var(--accent); padding-left: 16px; display: grid; gap: 4px; }
  h4 { font-size: 1.1rem; }
  code { font-family: ui-monospace, monospace; font-size: .75rem; color: var(--mute); margin-left: 8px; }
  .ref { border-top: 4px solid var(--rule); padding: 24px 0 96px; }
  h2 { font-size: 2rem; letter-spacing: -0.02em; }
  table { border-collapse: collapse; width: 100%; }
  tr { border-top: 1px solid var(--ink); }
  th, td { text-align: left; vertical-align: top; padding: 16px 16px 16px 0; }
  th { width: 25%; } th code { margin: 0; }
  .mute { color: var(--mute); }
  @media (max-width: 760px) {
    .grid { grid-template-columns: 1fr; gap: 24px; }
    .steps { grid-auto-flow: column; }
    .steps li { font-size: .75rem; flex-direction: column; align-items: flex-start; gap: 4px; }
    dl { grid-template-columns: 1fr; gap: 2px; } dd { margin-bottom: 12px; }
    th { width: auto; display: block; padding-bottom: 0; } td { display: block; }
  }
</style>
