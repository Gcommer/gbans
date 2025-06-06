BEGIN;


DO
$$
    BEGIN
        CREATE TYPE config_action AS ENUM (
            'ban',
            'gag',
            'kick'
            );
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;


ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_enabled boolean not null DEFAULT false;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_action config_action not null DEFAULT 'ban';

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_duration int not null DEFAULT 0;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_aim_snap int not null DEFAULT 40;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_psilent int not null DEFAULT 25;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_bhop int not null DEFAULT 20;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_fake_ang int not null DEFAULT 15;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_cmd_num int not null DEFAULT 40;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_too_many_connections int not null DEFAULT 1;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_cheat_cvar int not null DEFAULT 1;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_oob_var int not null DEFAULT 1;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS anticheat_max_invalid_user_cmd int not null DEFAULT 1;

ALTER TABLE config
    ADD COLUMN IF NOT EXISTS discord_anticheat_channel_id text not null DEFAULT '';

ALTER TABLE config
    DROP COLUMN IF EXISTS discord_unregister_on_start;

ALTER TABLE config
    DROP COLUMN IF EXISTS sentry_sentry_dsn;

ALTER TABLE config
    DROP COLUMN IF EXISTS sentry_sentry_dsn_web;

ALTER TABLE config
    DROP COLUMN IF EXISTS sentry_sentry_sample_rate;

ALTER TABLE config
    DROP COLUMN IF EXISTS sentry_sentry_trace;

COMMIT;
