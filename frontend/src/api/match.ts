import { PlayerClass, Team } from './const';
import { TeamScores } from './stats';
import { apiCall, QueryFilter } from './common';

interface MatchHealer {
    match_medic_id: number;
    match_player_id: number;
    healing: number;
    charges_uber: number;
    charges_kritz: number;
    charges_vacc: number;
    charges_quickfix: number;
    drops: number;
    near_full_charge_death: number;
    avg_uber_length: number;
    major_adv_lost: number;
    biggest_adv_lost: number;
}

export interface MatchPlayerClass {
    match_player_class_id: number;
    match_player_id: number;
    player_class: PlayerClass;
    kills: number;
    assists: number;
    deaths: number;
    playtime: number;
    dominations: number;
    dominated: number;
    revenges: number;
    damage: number;
    damage_taken: number;
    healing_taken: number;
    captures: number;
    captures_blocked: number;
    building_destroyed: number;
}

export interface MatchPlayerKillstreak {
    match_killstreak_id: number;
    match_player_id: number;
    player_class: PlayerClass;
    killstreak: number;
    duration: number;
}

export interface MatchPlayer {
    match_player_id: number;
    steam_id: string;
    team: Team;
    name: string;
    avatar_hash: string;
    time_start: string;
    time_end: string;
    kills: number;
    assists: number;
    deaths: number;
    suicides: number;
    dominations: number;
    dominated: number;
    revenges: number;
    damage: number;
    damage_taken: number;
    healing_taken: number;
    health_packs: number;
    healing_packs: number;
    captures: number;
    captures_blocked: number;
    extinguishes: number;
    building_built: number;
    building_destroyed: number;
    backstabs: number;
    airshots: number;
    headshots: number;
    shots: number;
    hits: number;
    medic_stats?: MatchHealer;
    classes: MatchPlayerClass[];
    killstreaks: MatchPlayerKillstreak[];
}

export interface PersonMessages {}

export interface MatchResult {
    match_id: string;
    server_id: number;
    title: string;
    map_name: string;
    team_scores: TeamScores;
    time_start: Date;
    time_end: Date;
    players: MatchPlayer[];
    chat: PersonMessages[];
}

export interface MatchSummary {}

export const apiGetMatch = async (match_id: string) =>
    await apiCall<MatchResult>(`/api/log/${match_id}`, 'GET');

export interface MatchesQueryOpts extends QueryFilter<MatchSummary> {
    steam_id?: string;
    server_id?: number;
    map?: string;
    time_start?: Date;
    time_end?: Date;
}

export const apiGetMatches = async (opts: MatchesQueryOpts) =>
    await apiCall<MatchSummary[], MatchesQueryOpts>(`/api/logs`, 'POST', opts);