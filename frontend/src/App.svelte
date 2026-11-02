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
  let size = $state('');

  try { size = localStorage.getItem('size') || ''; } catch {}
  $effect(() => {
    document.documentElement.className = size;
    try { localStorage.setItem('size', size); } catch {}
  });

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

<div class="flag" aria-hidden="true"><i></i><i></i><i></i></div>

<div class="util">
  <div class="wrap">
    <a class="skip" href="#main">{$t('skip')}</a>
    <div class="tools">
      <span class="sizes" role="group" aria-label={$t('textSize')}>
        <button onclick={() => (size = 't-sm')} aria-label="A-" class:on={size === 't-sm'}>A-</button>
        <button onclick={() => (size = '')} aria-label="A" class:on={size === ''}>A</button>
        <button onclick={() => (size = 't-lg')} aria-label="A+" class:on={size === 't-lg'}>A+</button>
        <button onclick={() => (size = 't-xl')} aria-label="A++" class:on={size === 't-xl'}>A++</button>
      </span>
      <label>
        <span class="sr">Language</span>
        <select bind:value={$lang}>
          {#each Object.entries(languages) as [code, name]}<option value={code}>{name}</option>{/each}
        </select>
      </label>
    </div>
  </div>
</div>

<header class="brand">
  <div class="wrap">
    <img src="/favicon.svg" alt="" width="56" height="56" />
    <div>
      <div class="name">{$t('title')}</div>
      <div class="tag">{$t('sub')}</div>
    </div>
  </div>
</header>

<nav aria-label="Main">
  <div class="wrap">
    <a href="#main" class="cur">{$t('navApply')}</a>
    <a href="#reasons">{$t('navReasons')}</a>
  </div>
</nav>

<div class="notice"><div class="wrap">ⓘ {$t('disclaimer')}</div></div>

<main id="main" class="wrap">
  <h1>{$t('pageTitle')}</h1>

  <div class="layout">
    <section class="panel" aria-labelledby="stepTitle">
      <ol class="steps" aria-label="Progress">
        {#each [1, 2, 3] as n}
          <li class:on={step === n} class:done={step > n} aria-current={step === n ? 'step' : undefined}>
            <span>{step > n ? '✓' : n}</span>{$t('step' + n)}
          </li>
        {/each}
      </ol>

      <h2 id="stepTitle">{$t('step' + step)}</h2>
      <p class="help">{$t('h' + step)}</p>

      <form onsubmit={(e) => e.preventDefault()}>
        {#if step === 1}
          <label>{$t('name')}<input bind:value={form.name} autocomplete="name" /></label>
          <label>{$t('dob')}<input bind:value={form.dob} inputmode="numeric" placeholder="01-01-1990" /></label>
          <label>{$t('mobile')}<input bind:value={form.mobile} inputmode="tel" maxlength="10" autocomplete="tel-national" /></label>
          <div class="row end"><button class="btn primary" onclick={() => go(2)}>{$t('next')} »</button></div>
        {:else if step === 2}
          <button class="btn detect" onclick={detect} disabled={detecting}>{detecting ? $t('detecting') : $t('detect')}</button>
          <p class="hint">{$t('detectHint')}</p>
          {#if locError}<p class="err" role="alert">{locError}</p>{/if}
          <div class="two">
            <label>{$t('village')}<input bind:value={form.address.village} /></label>
            <label>{$t('district')}<input bind:value={form.address.district} /></label>
            <label>{$t('state')}<input bind:value={form.address.state} /></label>
            <label>{$t('pincode')}<input bind:value={form.pincode} inputmode="numeric" maxlength="6" /></label>
          </div>
          <div class="row">
            <button class="btn" onclick={() => go(1)}>« {$t('back')}</button>
            <button class="btn primary" onclick={() => go(3)}>{$t('next')} »</button>
          </div>
        {:else}
          <dl>
            <dt>{$t('name')}</dt><dd>{form.name || '—'}</dd>
            <dt>{$t('dob')}</dt><dd>{form.dob || '—'}</dd>
            <dt>{$t('mobile')}</dt><dd>{form.mobile || '—'}</dd>
            <dt>{$t('step2')}</dt><dd>{[form.address.village, form.address.district, form.address.state, form.pincode].filter(Boolean).join(', ') || '—'}</dd>
          </dl>
          <div class="row">
            <button class="btn" onclick={() => go(2)}>« {$t('back')}</button>
            <button class="btn primary" onclick={submit}>{$t('submit')}</button>
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
                    <p><b>{$t('why')}:</b> {i.why}</p>
                    <p><b>{$t('fix')}:</b> {i.fix}</p>
                  </article>
                {/each}
              {/if}
            </div>
          {/if}
        {/if}
      </form>
    </section>

    <aside>
      <h2>{$t('helpline')}</h2>
      <p class="big">14555</p>
    </aside>
  </div>

  <section id="reasons" class="reasons">
    <h2>{$t('reasons')}</h2>
    <p class="help">{$t('reasonsSub')}</p>
    <div class="scroll">
      <table>
        <thead><tr><th scope="col">{$t('colCode')}</th><th scope="col">{$t('colReason')}</th><th scope="col">{$t('fix')}</th></tr></thead>
        <tbody>
          {#each reasons as r}
            <tr><td><code>{r.code}</code></td><td><b>{r.title}</b><br />{r.why}</td><td>{r.fix}</td></tr>
          {/each}
        </tbody>
      </table>
    </div>
  </section>
</main>

<footer><div class="wrap">{$t('disclaimer')}</div></footer>

<style>
  .wrap { max-width: 1100px; margin: 0 auto; padding: 0 20px; }
  .sr { position: absolute; left: -999px; }
  .flag { display: flex; height: 6px; } .flag i { flex: 1; } .flag i:nth-child(1) { background: var(--saffron); } .flag i:nth-child(2) { background: #fff; box-shadow: inset 0 1px #ddd, inset 0 -1px #ddd; } .flag i:nth-child(3) { background: var(--green); }

  .util { background: #f1f3f5; border-bottom: 1px solid var(--line); font-size: .8rem; }
  .util .wrap { display: flex; justify-content: space-between; align-items: center; gap: 12px; min-height: 40px; flex-wrap: wrap; }
  .skip { color: var(--ink); }
  .tools { display: flex; align-items: center; gap: 12px; }
  .sizes { display: flex; }
  .sizes button { border: 1px solid var(--line); background: #fff; padding: 2px 8px; cursor: pointer; margin-left: -1px; min-width: 36px; }
  .sizes button.on { background: var(--navy); color: #fff; border-color: var(--navy); }
  select { border: 1px solid #888; background: #fff; padding: 3px 6px; }

  .brand { border-bottom: 1px solid var(--line); }
  .brand .wrap { display: flex; align-items: center; gap: 16px; padding-top: 14px; padding-bottom: 14px; }
  .name { font-size: 1.8rem; font-weight: 700; color: var(--navy); line-height: 1.1; }
  .tag { color: var(--mute); font-size: .95rem; }

  nav { background: var(--navy); }
  nav .wrap { display: flex; gap: 0; flex-wrap: wrap; }
  nav a { color: #fff; text-decoration: none; padding: 12px 20px; font-weight: 600; border-right: 1px solid rgba(255,255,255,.2); }
  nav a:hover, nav a.cur { background: var(--blue); }

  .notice { background: var(--warn); border-bottom: 1px solid #e6cf8b; font-size: .85rem; }
  .notice .wrap { padding-top: 8px; padding-bottom: 8px; }

  main { padding-top: 24px; padding-bottom: 56px; }
  h1 { font-size: 1.6rem; color: var(--navy); border-bottom: 3px solid var(--saffron); display: inline-block; padding-bottom: 4px; margin-bottom: 20px; }
  h2 { font-size: 1.2rem; color: var(--navy); }
  .help { color: var(--mute); margin: 2px 0 16px; }

  .layout { display: grid; grid-template-columns: 1fr 260px; gap: 24px; align-items: start; }
  .panel { border: 1px solid var(--line); border-top: 4px solid var(--navy); padding: 20px 24px 24px; background: #fff; }
  aside { background: var(--pale); border: 1px solid var(--line); padding: 16px 20px; }
  .big { font-size: 2.2rem; font-weight: 700; color: var(--navy); }

  .steps { list-style: none; padding: 0; display: flex; gap: 0; margin-bottom: 20px; border: 1px solid var(--line); }
  .steps li { flex: 1; padding: 8px 12px; display: flex; gap: 8px; align-items: center; font-size: .85rem; color: var(--mute); background: var(--pale); border-right: 1px solid var(--line); }
  .steps li:last-child { border-right: 0; }
  .steps span { width: 24px; height: 24px; border-radius: 50%; border: 1px solid currentColor; display: grid; place-items: center; flex: none; font-size: .8rem; }
  .steps li.on { background: var(--navy); color: #fff; font-weight: 600; }
  .steps li.done { color: var(--green); font-weight: 600; }

  form { display: grid; gap: 16px; }
  label { display: grid; gap: 4px; font-weight: 600; font-size: .9rem; }
  input { border: 1px solid #777; padding: 9px 10px; background: #fff; font-size: 1rem; border-radius: 2px; }
  input:focus { outline: 3px solid var(--saffron); border-color: var(--navy); }
  .two { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
  .row { display: flex; justify-content: space-between; gap: 12px; flex-wrap: wrap; }
  .row.end { justify-content: flex-end; }

  .btn { border: 1px solid var(--navy); background: #fff; color: var(--navy); padding: 9px 22px; font-weight: 600; cursor: pointer; border-radius: 2px; min-height: 44px; }
  .btn:hover:not(:disabled) { background: var(--pale); }
  .btn.primary { background: var(--navy); color: #fff; }
  .btn.primary:hover:not(:disabled) { background: var(--blue); }
  .btn.detect { background: var(--saffron); border-color: #c46f00; color: #1a1a1a; width: 100%; }
  .btn:disabled { opacity: .6; cursor: wait; }
  .hint { color: var(--mute); font-size: .85rem; margin-top: -8px; }
  .err { color: var(--err); font-weight: 600; }

  dl { display: grid; grid-template-columns: 1fr 2fr; gap: 8px 20px; border: 1px solid var(--line); padding: 14px 16px; background: var(--pale); }
  dt { color: var(--mute); font-size: .85rem; } dd { font-weight: 600; }
  .result { display: grid; gap: 12px; border-top: 1px solid var(--line); padding-top: 16px; }
  .ok { color: var(--green); font-weight: 700; font-size: 1.15rem; border: 1px solid var(--green); background: #eaf6e8; padding: 12px 16px; }
  .result article { border: 1px solid #e3a8a3; border-left: 5px solid var(--err); background: #fdf1f0; padding: 10px 14px; }
  h3 { font-size: 1.05rem; color: var(--err); }
  h4 { font-size: 1rem; }
  code { font-family: ui-monospace, Menlo, monospace; font-size: .75rem; background: var(--pale); border: 1px solid var(--line); padding: 1px 6px; margin-left: 4px; }
  td code { margin: 0; }

  .reasons { margin-top: 48px; }
  .scroll { overflow-x: auto; }
  table { border-collapse: collapse; width: 100%; font-size: .9rem; }
  th { background: var(--navy); color: #fff; text-align: left; padding: 10px 12px; }
  td { padding: 10px 12px; vertical-align: top; border: 1px solid var(--line); }
  tbody tr:nth-child(even) { background: var(--pale); }

  footer { background: var(--navy); color: #dbe6f2; font-size: .85rem; padding: 18px 0; }

  @media (max-width: 760px) {
    .layout { grid-template-columns: 1fr; }
    .two { grid-template-columns: 1fr; }
    .steps { flex-direction: column; } .steps li { border-right: 0; border-bottom: 1px solid var(--line); }
    dl { grid-template-columns: 1fr; gap: 2px; } dd { margin-bottom: 8px; }
    .panel { padding: 16px; }
    nav a { flex: 1; text-align: center; padding: 10px; }
  }
</style>
