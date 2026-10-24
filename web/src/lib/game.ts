export type Player = {
	name: string;
};

export type Game = {
	id: string;
	players: Player[];
	version: number;
	prompt: string;
};
