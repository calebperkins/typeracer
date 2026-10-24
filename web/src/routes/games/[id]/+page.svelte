<script lang="ts">
	import type { PageProps } from './$types';
	import { resolve } from '$app/paths';

	import { onMount } from 'svelte';

	let { data }: PageProps = $props();

	let version = $state<string | null>(null);

	onMount(() => {
		const proto = location.protocol === 'https:' ? 'wss' : 'ws';
		const ws = new WebSocket(`${proto}://${location.host}/api/games/${data.game.id}/ws`);
		ws.onmessage = (e) => {
			version = JSON.parse(e.data).data;
		};
		return () => ws.close();
	});

	const refresh = () => fetch(`/api/games/${data.game.id}/`, { method: 'POST' });
</script>

<a href={resolve('/')}>&larr; Games</a>

<h1>Game {data.game.id}</h1>

{#if data.game.players.length === 0}
	<p>No players yet.</p>
{:else}
	<ul>
		{#each data.game.players as player (player.name)}
			<li>{player.name}</li>
		{/each}
	</ul>
{/if}


<h2>Debug</h2>
<p>Version: {version ?? 'connecting…'}</p>
<button onclick={refresh}>Refresh</button>
