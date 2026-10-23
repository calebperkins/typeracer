<script lang="ts">
	import type { PageProps } from './$types';
	import { resolve } from '$app/paths';

	let { data }: PageProps = $props();
</script>

<h1>Games</h1>

{#if data.games.length === 0}
	<p>No games yet.</p>
{:else}
	<ul>
		{#each data.games as game (game.id)}
			<li><a href={resolve('/games/[id]', { id: game.id })}>{game.id}</a></li>
		{/each}
	</ul>
{/if}
