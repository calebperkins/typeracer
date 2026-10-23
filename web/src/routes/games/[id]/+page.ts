import { error } from '@sveltejs/kit';
import type { Game } from '#lib/game.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`/api/games/${encodeURIComponent(params.id)}`);
	if (!res.ok) {
		error(res.status, 'Failed to load game');
	}
	const game: Game = await res.json();
	return { game };
};
