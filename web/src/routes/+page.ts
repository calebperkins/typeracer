import { error } from '@sveltejs/kit';
import type { Game } from '#lib/game.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch('/api/games');
	if (!res.ok) {
		error(res.status, 'Failed to load games');
	}
	const games: Game[] = await res.json();
	return { games };
};
