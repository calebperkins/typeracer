<script lang="ts">
	import type { PageProps } from './$types';
	import { resolve } from '$app/paths';

	let { data }: PageProps = $props();
</script>

<a href={resolve('/')}>&larr; Games</a>

<h1>Game {data.game.id}</h1>
